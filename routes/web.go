package routes

import (
	handlers "github.com/deeincom/deeincom/app/handlers/web"

	"github.com/gofiber/fiber/v2"
)

func RegisterWeb(web fiber.Router) {
	h := handlers.NewHandler()
	// Homepage
	web.Get("/", h.Index)
}
