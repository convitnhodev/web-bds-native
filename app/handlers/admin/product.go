package web

import (
	"context"
	"github.com/kurin/blazer/b2"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"path/filepath"
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

func (h *handler) ProductUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	id := "002026794d2122d000000000c"
	key := "K0022S/oKffHhTFf12PpuVQjim2O0Xw"

	ctx := context.Background()
	b2, err := b2.NewClient(ctx, id, key)
	if err != nil {
		log.Fatalln(err)
	}
	bucket, _ := b2.Bucket(ctx, "d1e2e3i4n5")
	obj := bucket.Object(filepath.Join("tmp", file.Filename))
	w := obj.NewWriter(ctx)
	defer w.Close()
	if _, err := io.Copy(w, src); err != nil {
		return err
	}
	return c.JSON(200, echo.Map{"status": "ok", "object_name": obj.Name()})
}
