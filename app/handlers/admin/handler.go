package web

import (
	"github.com/deeincom/deeincom/app/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
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
	return c.Render(http.StatusOK, "admin.dashboard.page.html", echo.Map{
		"Context": c,
	})
}

func (h *handler) UserList(c echo.Context) error {
	type User struct {
		ID          int64  `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}
	return c.Render(http.StatusOK, "admin.user.list.page.html", echo.Map{
		"Context": c,
		"Users": []User{
			{
				ID:          1,
				FirstName:   "Huy",
				LastName:    "Huynh",
				Email:       "huyhvq@icloud.com",
				PhoneNumber: "+84123123123",
			},
		},
	})
}
