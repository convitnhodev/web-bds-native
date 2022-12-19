package web

import (
	"context"
	"github.com/deeincom/deeincom/app/models"
	"github.com/kurin/blazer/b2"
	"github.com/labstack/echo-contrib/session"
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
	Title                string `json:"title" form:"title" validate:"required"`
	CategoryID           string `json:"category_id" form:"category_id" validate:"required"`
	Description          string `json:"description" form:"description" validate:"required"`
	City                 string `json:"city" form:"city" validate:"required"`
	District             string `json:"district" form:"district" validate:"required"`
	Ward                 string `json:"ward" form:"ward" validate:"required"`
	AddressNumber        string `json:"address_number" form:"address_number" validate:"required"`
	Street               string `json:"street" form:"street" validate:"required"`
	Area                 int64  `json:"area" form:"area" validate:"required"`
	DocumentType         string `json:"document_type" form:"document_type" validate:"required"`
	DocumentProof        string `json:"document_proof" form:"document_proof" validate:"required"`
	Bedroom              int64  `json:"bedroom" form:"bedroom" validate:"required"`
	Toilet               int64  `json:"toilet" form:"toilet" validate:"required"`
	Floor                int64  `json:"floor" form:"floor" validate:"required"`
	HouseDirection       string `json:"house_direction" form:"house_direction" validate:"required"`
	BalconyDirection     string `json:"balcony_direction" form:"balcony_direction" validate:"required"`
	FrontWidth           int64  `json:"front_width" form:"front_width" validate:"required"`
	StreetWidth          int64  `json:"street_width" form:"street_width" validate:"required"`
	PavementWidth        int64  `json:"pavement_width" form:"pavement_width" validate:"required"`
	BusinessAdvantage    string `json:"business_advantage" form:"business_advantage" validate:"required"`
	FinancialPlan        string `json:"financial_plan" form:"financial_plan" validate:"required"`
	GeneralPurchaseTerms string `json:"general_purchase_terms" form:"general_purchase_terms" validate:"required"`
	Furniture            string `json:"furniture" form:"furniture" validate:"required"`
}

func (h *handler) ProductStore(c echo.Context) error {
	var p Product
	if err := c.Bind(&p); err != nil {
		sess, _ := session.Get("error_message", c)
		sess.AddFlash(err.Error(), "message")
		_ = sess.Save(c.Request(), c.Response())
		return c.Redirect(301, "/admin/products/create")
	}
	if err := h.validator.Struct(&p); err != nil {
		sess, _ := session.Get("error_message", c)
		sess.AddFlash(err.Error(), "message")
		_ = sess.Save(c.Request(), c.Response())
		return c.Redirect(301, "/admin/products/create")
	}

	if err := h.repository.Product.Create(models.Product{
		Title:                p.Title,
		CategoryID:           p.CategoryID,
		Description:          p.Description,
		City:                 p.City,
		District:             p.District,
		Ward:                 p.Ward,
		AddressNumber:        p.AddressNumber,
		Street:               p.Street,
		Area:                 p.Area,
		DocumentType:         p.DocumentType,
		DocumentProof:        p.DocumentProof,
		Bedroom:              p.Bedroom,
		Toilet:               p.Toilet,
		Floor:                p.Floor,
		HouseDirection:       p.HouseDirection,
		BalconyDirection:     p.BalconyDirection,
		FrontWidth:           p.FrontWidth,
		StreetWidth:          p.StreetWidth,
		PavementWidth:        p.PavementWidth,
		BusinessAdvantage:    p.BusinessAdvantage,
		FinancialPlan:        p.FinancialPlan,
		Furniture:            p.Furniture,
		GeneralPurchaseTerms: p.GeneralPurchaseTerms,
		IsActivated:          true,
	}); err != nil {
		sess, _ := session.Get("error_message", c)
		sess.AddFlash(err.Error(), "message")
		_ = sess.Save(c.Request(), c.Response())
		return c.Redirect(301, "/admin/products/create")
	}
	sess, _ := session.Get("success_message", c)
	sess.AddFlash("tạo sản phẩm thành công", "message")
	_ = sess.Save(c.Request(), c.Response())
	return c.Redirect(301, "/admin/products/create")

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
