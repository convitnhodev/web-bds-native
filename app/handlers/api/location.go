package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type locationAPI handler

func (h *locationAPI) ListProvinces(c echo.Context) error {
	provinces, err := h.api.repository.Location.ListProvinces()
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, provinces)
}

func (h *locationAPI) ListDistrictsByProvinceID(c echo.Context) error {
	fmt.Println("ID:", c.Param("id"))
	r, err := h.api.repository.Location.ListDistrictsByProvinceID(c.Param("id"))
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, r)
}

func (h *locationAPI) ListWardsByDistrictID(c echo.Context) error {
	fmt.Println("ID:", c.Param("id"))
	r, err := h.api.repository.Location.ListWardsByDistrictID(c.Param("id"))
	if err != nil {
		return errJson(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, r)
}
