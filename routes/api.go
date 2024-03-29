package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/api"
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/jwt"
	"github.com/deeincom/deeincom/pkg/sms"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterAPI(web *echo.Echo, r *repositories.Repository, cfg *config.Config) {
	jwtAuth := jwt.NewAuth(cfg.GetString("AUTH_SECRET"), cfg.GetDuration("AUTH_EXPIRE_DURATION"))
	smsClient := sms.NewSMSClient(
		cfg.GetString("SMS_KEY"),
		cfg.GetString("SMS_SECRET_KEY"),
		cfg.GetString("SMS_BRAND_NAME"),
	)
	h := handlers.NewHandler(r, jwtAuth, smsClient)
	api := web.Group("/api/v1")
	api.POST("/accounts", h.AccountHandler.Register)
	api.POST("/accounts/auth", h.AccountHandler.Auth)

	api.GET("/locations/provinces", h.LocationHandler.ListProvinces)
	api.GET("/locations/provinces/:id/districts", h.LocationHandler.ListDistrictsByProvinceID)
	api.GET("/locations/districts/:id/wards", h.LocationHandler.ListWardsByDistrictID)

	api.Use(middleware.JWTWithConfig(jwtAuth.GetMiddlewareConfig()))
	api.GET("/accounts/profile", h.AccountHandler.Profile)
	api.POST("/accounts/verification", h.AccountHandler.SentCode)
	api.POST("/accounts/verification/:code/code", h.AccountHandler.VerifyCode)
}
