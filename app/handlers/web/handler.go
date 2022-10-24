package web

import (
	"net/http"

	"github.com/deeincom/deeincom/app/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type handler struct {
	validator  *validator.Validate
	repository *repositories.Repository
}

func NewHandler(r *repositories.Repository) *handler {
	var validate = validator.New()
	return &handler{
		validator:  validate,
		repository: r,
	}
}

func (h *handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "home.page.html", map[string]string{
		"Title": "Hello, World!",
	})
}

func (h *handler) Detail(c echo.Context) error {
	return c.Render(http.StatusOK, "detail.page.html", map[string]string{
		"Title": "Hello, World!",
	})
}

func (h *handler) Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register.page.html", map[string]string{
		"Title": "Hello, World!",
	})
}

func (h *handler) Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.page.html", map[string]string{
		"Title": "Hello, World!",
	})
}

func (h *handler) Verify(c echo.Context) error {
	return c.Render(http.StatusOK, "verify.page.html", map[string]string{
		"Title": "Hello, World!",
	})
}
