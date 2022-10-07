package web

import (
	"github.com/gofiber/fiber/v2"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}
