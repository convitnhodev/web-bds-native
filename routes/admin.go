package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/admin"
	"github.com/labstack/echo/v4"
)

func RegisterAdmin(web *echo.Echo) {
	h := handlers.NewHandler()
	// Homepage
	r := web.Group("/admin")
	r.GET("", h.Index)
	r.GET("/users", h.UserList)
	r.GET("/products/create", h.ProductCreate)
	r.POST("/products", h.ProductCreate)
	r.POST("/products/images/upload", h.ProductUpload)
}
