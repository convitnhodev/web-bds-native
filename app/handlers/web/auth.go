package web

import (
	"github.com/labstack/echo/v4"
)

type RegisterCreateRequest struct {
	FirstName            string `form:"first_name" validate:"required"`
	LastName             string `form:"last_name" validate:"required"`
	Email                string `form:"email" validate:"email"`
	Password             string `form:"password" validate:"required"`
	PasswordConfirmation string `form:"password_confirmation" validate:"required,eqfield=Password"`
	PhoneNumber          string `form:"phone_number" validate:"required,e164"`
}

func (h *handler) RegisterCreate(c echo.Context) error {
	req := RegisterCreateRequest{}
	if err := c.Bind(&req); err != nil {
		//return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		return err
	}
	if err := h.validator.Struct(&req); err != nil {
		//return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		return err
	}
	pe, err := h.repository.User.IsPhoneExist(req.PhoneNumber)
	if err != nil {
		//return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		return err
	}
	if pe {
		//return fiber.NewError(fiber.StatusUnprocessableEntity, "phone number existed")
		return nil
	}
	//return fiber.NewError(fiber.StatusOK, "Created")
	return nil
}
