package web

import (
	"context"
	"fmt"
	"github.com/kurin/blazer/b2"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"os"
)

func (h *handler) ProductCreate(c echo.Context) error {
	accessKeyID, policy, signature := GetS3BrowserSecret()

	//id := "002026794d2122d000000000c"
	//key := "K0022S/oKffHhTFf12PpuVQjim2O0Xw"
	//
	//ctx := context.Background()
	//
	//// b2_authorize_account
	//b2, err := b2.NewClient(ctx, id, key)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//buckets, err := b2.ListBuckets(ctx)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//for _, bucket := range buckets {
	//	fmt.Println(bucket.Name())
	//}

	return c.Render(http.StatusOK, "admin.product.create.page.html", echo.Map{
		"Context":       c,
		"Furniture":     furniture,
		"S3AccessKeyID": accessKeyID,
		"S3Policy":      policy,
		"S3Signature":   signature,
		"S3Endpoint":    s3Endpoint,
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

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	id := "002026794d2122d000000000c"
	key := "K0022S/oKffHhTFf12PpuVQjim2O0Xw"

	ctx := context.Background()

	// b2_authorize_account
	b2, err := b2.NewClient(ctx, id, key)
	if err != nil {
		log.Fatalln(err)
	}

	buckets, err := b2.ListBuckets(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("buckets: %+v\n", buckets)

	bucket, _ := b2.Bucket(ctx, "d1e2e3i4n5")
	obj := bucket.Object("test/test-image.jpg")
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, src); err != nil {
		w.Close()
		return err
	}
	w.Close()

	return c.JSON(200, echo.Map{"status": "ok"})
}
