package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard.page.html", nil)
}

func (h *handler) UserList(c echo.Context) error {
	type User struct {
		ID          int64  `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}
	return c.Render(http.StatusOK, "user.list.page.html", map[string]interface{}{
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
