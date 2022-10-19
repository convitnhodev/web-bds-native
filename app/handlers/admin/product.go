package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *handler) ProductCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "admin.product.create.page.html", echo.Map{
		"Context": c,
	})
}

type Product struct {
	Title             string `json:"title" form:"title"`
	CategoryID        string `json:"category_id" form:"category_id"`
	Description       string `json:"description" form:"description"`
	City              string `json:"city" form:"city"`
	District          string `json:"district" form:"district"`
	Ward              string `json:"ward" form:"ward"`
	AddressNumber     string `json:"address_number" form:"address_number"`
	Street            string `json:"street" form:"street"`
	Area              int64  `json:"area" form:"area"`
	DocumentType      string `json:"document_type" form:"document_type"`
	DocumentProof     string `json:"document_proof" form:"document_proof"`
	Bedroom           int64  `json:"bedroom" form:"bedroom"`
	Toilet            int64  `json:"toilet" form:"toilet"`
	Floor             int64  `json:"floor" form:"floor"`
	HouseDirection    string `json:"house_direction" form:"house_direction"`
	BalconyDirection  string `json:"balcony_direction" form:"balcony_direction"`
	FrontWidth        int64  `json:"front_width" form:"front_width"`
	StreetWidth       int64  `json:"street_width" form:"street_width"`
	PavementWidth     int64  `json:"pavement_width" form:"pavement_width"`
	BusinessAdvantage string `json:"business_advantage" form:"business_advantage"`
	FinancialPlan     string `json:"financial_plan" form:"financial_plan"`
	Furniture         string `json:"furniture" form:"furniture"`
}

func (h *handler) ProductStore(c echo.Context) error {
	return nil
}
