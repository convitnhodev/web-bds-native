package web

import (
	"fmt"
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
	products, err := h.repository.Product.ListProducts()
	if err != nil {
		//handle error
	}
	wards := make([]string, len(products))
	for _, product := range products {
		wards = append(wards, product.Ward)
	}

	locations, err := h.repository.Location.ListLocationByWardIDs(wards)
	if err != nil {
		//handle error
	}
	for i, product := range products {
		if l, ok := locations[product.Ward]; ok {
			products[i].Ward = l.WardName
			products[i].District = l.DistrictName
			products[i].City = l.ProvinceName
		}
	}

	fmt.Printf("products: %+v\n", products)

	return c.Render(http.StatusOK, "home.page.html", echo.Map{
		"Title":    "Hello, World!",
		"Products": products,
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
