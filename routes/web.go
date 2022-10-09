package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/web"
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/labstack/echo/v4"
)

func RegisterWeb(web *echo.Echo, r *repositories.Repository) {
	h := handlers.NewHandler(r)
	// Homepage
	web.GET("/", h.Index)
}
