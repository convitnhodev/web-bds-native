package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/api"
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/deeincom/deeincom/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func RegisterAPI(web *echo.Echo, r *repositories.Repository, jwt jwt.Authenticator) {
	h := handlers.NewHandler(r, jwt)

	api := web.Group("/api/v1")
	api.POST("/accounts", h.UserHandler.Register)
	api.POST("/accounts/auth", h.UserHandler.Auth)
}
