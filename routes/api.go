package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/api"
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func RegisterAPI(web *echo.Echo, r *repositories.Repository, cfg *config.Config) {
	jwtAuth := jwt.NewAuth(cfg.GetString("AUTH_SECRET"), cfg.GetDuration("AUTH_EXPIRE_DURATION"))
	h := handlers.NewHandler(r, jwtAuth)
	api := web.Group("/api/v1")
	//api.Use(middleware.JWTWithConfig(jwtAuth.GetMiddlewareConfig()))
	api.GET("/accounts/profile", h.AccountHandler.Profile)
	api.POST("/accounts", h.AccountHandler.Register)
	api.POST("/accounts/auth", h.AccountHandler.Auth)
	api.POST("/accounts/verification", h.AccountHandler.SentCode)
	api.POST("/accounts/verification/:code/code", h.AccountHandler.VerifyCode)

	api.GET("/locations/provinces", h.LocationHandler.ListProvinces)
	api.GET("/locations/provinces/:id/districts", h.LocationHandler.ListDistrictsByProvinceID)
	api.GET("/locations/districts/:id/wards", h.LocationHandler.ListWardsByDistrictID)
}
