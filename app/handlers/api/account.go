package api

import (
	"errors"
	"fmt"
	"github.com/deeincom/deeincom/app/models"
	authJwt "github.com/deeincom/deeincom/pkg/jwt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type accountAPI handler

type registerReq struct {
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	PhoneNumber          string `json:"phone_number" validate:"required,e164"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type loginReq struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Password    string `json:"password" validate:"required"`
}

type accountResp struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	IsVerified  bool      `json:"is_verified"`
	IsActivated bool      `json:"is_activated"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type authResp struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// Register godoc
// @Summary      register
// @Description  register an account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param data body registerReq true "Register account information"
// @Success      200  {object}  accountResp
// @Failure      422  {object}  defaultJsonResp
// @Failure      500  {object}  defaultJsonResp
// @Router       /accounts [post]
func (h *accountAPI) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	if err := h.api.validator.Struct(&req); err != nil {
		return errJson(c, http.StatusUnprocessableEntity, err)
	}
	user, err := h.api.repository.User.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil && err != db.ErrNoMoreRows {
		return errJson(c, http.StatusInternalServerError, err)
	}
	if user != nil {
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	mu := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hash),
		IsActivated: true,
	}
	user, err = h.api.repository.User.CreateAnUser(&mu)
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, &accountResp{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  isAccountVerified(user.VerifiedAt),
		IsActivated: user.IsActivated,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

// Auth godoc
// @Summary      authenticate
// @Description  get token using credentials
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param data body loginReq true "account credentials"
// @Success      200  {object}  authResp
// @Failure      422  {object}  defaultJsonResp
// @Failure      500  {object}  defaultJsonResp
// @Router       /accounts/auth [post]
func (h *accountAPI) Auth(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	if err := h.api.validator.Struct(&req); err != nil {
		return errJson(c, http.StatusUnprocessableEntity, err)
	}
	user, err := h.api.repository.User.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	jwtToken, err := h.api.jwt.Issue(user.PhoneNumber, isAccountVerified(user.VerifiedAt))
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}

	sess, _ := session.Get("auth", c)
	sess.Values["user_id"] = user.ID
	sess.Values["user_is_verified"] = isAccountVerified(user.VerifiedAt)
	_ = sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, &authResp{
		Token: jwtToken,
		Type:  "Bearer",
	})
}

// Profile godoc
// @Summary      account profiles
// @Description  get account profile
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success		200  {object}  accountResp
// @Failure     422  {object}  defaultJsonResp
// @Failure		500  {object}  defaultJsonResp
// @Security	ApiBearerKey
// @Router		/accounts/profile [get]
func (h *accountAPI) Profile(c echo.Context) error {
	claims := getClaims(c)
	user, err := h.api.repository.User.GetUserByPhoneNumber(claims.PhoneNumber)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	fmt.Printf("USER: %+v\n", user)
	return c.JSON(http.StatusOK, &accountResp{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  isAccountVerified(user.VerifiedAt),
		IsActivated: user.IsActivated,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

// SentCode godoc
// @Summary      sent verify code
// @Description  sent sms verify code
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success		200  {object}  defaultJsonResp
// @Failure     422  {object}  defaultJsonResp
// @Failure		500  {object}  defaultJsonResp
// @Security	ApiBearerKey
// @Router		/accounts/verification [get]
func (h *accountAPI) SentCode(c echo.Context) error {
	claims := getClaims(c)
	user, err := h.api.repository.User.GetUserByPhoneNumber(claims.PhoneNumber)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	rn := rand.Intn(max-min+1) + min
	if err := h.api.repository.User.CreateVerifyCode(user.PhoneNumber, strconv.Itoa(rn)); err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	// sent email or sms
	return successJson(c, http.StatusCreated, "sms has been send")
}

// VerifyCode godoc
// @Summary      sent verify code
// @Description  sent sms verify code
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        code   path      int  true  "verify code"
// @Success		200  {object}  authResp
// @Failure     422  {object}  defaultJsonResp
// @Failure		500  {object}  defaultJsonResp
// @Security	ApiBearerKey
// @Router		/accounts/verification/{code}/code [get]
func (h *accountAPI) VerifyCode(c echo.Context) error {
	token := c.Param("code")
	claims := getClaims(c)
	user, err := h.api.repository.User.GetUserByPhoneNumber(claims.PhoneNumber)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your phone number and password is incorrect"))
	}
	if isAccountVerified(user.VerifiedAt) {
		return errJson(c, http.StatusUnprocessableEntity, errors.New("this account has been verified"))
	}
	_, err = h.api.repository.User.GetVerifyCode(claims.PhoneNumber, token)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
		return errJson(c, http.StatusUnprocessableEntity, errors.New("your verify code has expired or invalid"))
	}
	if err := h.api.repository.User.UpdateVerified(claims.PhoneNumber); err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	jwtToken, err := h.api.jwt.Issue(user.PhoneNumber, true)
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	
	sess, _ := session.Get("auth", c)
	sess.Values["user_is_verified"] = isAccountVerified(user.VerifiedAt)
	_ = sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, &authResp{
		Token: jwtToken,
		Type:  "Bearer",
	})
}

func getClaims(c echo.Context) *authJwt.Claim {
	jwtToken := c.Get("user").(*jwt.Token)
	return jwtToken.Claims.(*authJwt.Claim)
}

func isAccountVerified(t time.Time) bool {
	fmt.Printf("is_zero %+v\n", t.IsZero())
	return !t.IsZero()
}
