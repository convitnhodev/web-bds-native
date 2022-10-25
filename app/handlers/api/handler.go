package api

import (
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/deeincom/deeincom/pkg/jwt"
	"github.com/deeincom/deeincom/pkg/sms"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type handler struct {
	api *Handler
}

type Handler struct {
	jwt        jwt.Authenticator
	validator  *validator.Validate
	repository *repositories.Repository
	common     handler
	smsSender  *sms.SMSSender

	AccountHandler  *accountAPI
	LocationHandler *locationAPI
}

type defaultJsonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHandler(r *repositories.Repository, jwt jwt.Authenticator, s *sms.SMSSender) *Handler {
	var validate = validator.New()
	h := &Handler{
		validator:  validate,
		repository: r,
		jwt:        jwt,
		smsSender:  s,
	}
	h.common.api = h
	h.AccountHandler = (*accountAPI)(&h.common)
	h.LocationHandler = (*locationAPI)(&h.common)
	return h
}

func errJson(c echo.Context, code int, err error) error {
	return c.JSON(code, &defaultJsonResp{
		Code:    code,
		Message: err.Error(),
	})
}

func successJson(c echo.Context, code int, msg string) error {
	return c.JSON(code, &defaultJsonResp{
		Code:    code,
		Message: msg,
	})
}
