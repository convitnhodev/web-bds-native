package api

import (
	"errors"
	"github.com/deeincom/deeincom/app/models"
	"github.com/labstack/echo/v4"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AccountHandler interface {
	Register(c echo.Context) error
	Auth(c echo.Context) error
}

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
	Email       *string   `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	VerifiedAt  time.Time `json:"verified_at,omitempty"`
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
// @Failure      422  {object}  errorResp
// @Failure      500  {object}  errorResp
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
	if err != nil {
		if err != db.ErrNoMoreRows {
			return errJson(c, http.StatusInternalServerError, err)
		}
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
		IsActivated: false,
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
		VerifiedAt:  user.VerifiedAt,
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
// @Failure      422  {object}  errorResp
// @Failure      500  {object}  errorResp
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
	isVerified := false
	if user.VerifiedAt.After(time.Date(2001, 2, 3, 4, 5, 6, 700000000, time.UTC)) {
		isVerified = true
	}
	jwtToken, err := h.api.jwt.Issue(user.PhoneNumber, isVerified)
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, &authResp{
		Token: jwtToken,
		Type:  "Bearer",
	})
}
