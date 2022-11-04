package models

import (
	"fmt"
	"time"

	"github.com/deeincom/deeincom/pkg/form"
)

type User struct {
	ID                  int
	Email               string
	Phone               string
	Password            string
	FirstName           string
	LastName            string
	Roles               []string
	EmailToken          string
	PhoneToken          string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	SendVerifiedEmailAt time.Time
	SendVerifiedPhoneAt time.Time
}

type Attachment struct {
	ID             int
	Title          string
	ContentType    string
	MineType       string
	Link           string
	Width          int
	Height         int
	Size           int
	VideoLength    int
	VideoThumbnail string
	Product        Product
}

type Product struct {
	ID                int
	Title             string
	Short             string
	Full              string
	FullContent       string
	City              string
	District          string
	Ward              string
	AddressNumber     string
	Street            string
	HouseDirection    string
	BalconyDirection  string
	BusinessAdvantage string
	FinancialPlan     string
	Legal             string
	Furniture         string
	Slug              string
	Type              string
	Area              int
	Bedroom           int
	Toilet            int
	Floor             int
	FrontWidth        int
	StreetWidth       int
	PavementWidth     int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (o *Product) Form() *form.Form {
	f := form.New(nil)
	f.Set("Title", o.Title)
	f.Set("Short", o.Short)
	f.Set("Full", o.Full)
	f.Set("FullContent", o.FullContent)
	f.Set("City", o.City)
	f.Set("District", o.District)
	f.Set("Ward", o.Ward)
	f.Set("AddressNumber", o.AddressNumber)
	f.Set("Street", o.Street)
	f.Set("HouseDirection", o.HouseDirection)
	f.Set("BalconyDirection", o.BalconyDirection)
	f.Set("BusinessAdvantage", o.BusinessAdvantage)
	f.Set("FinancialPlan", o.FinancialPlan)
	f.Set("Furniture", o.Furniture)
	f.Set("Type", o.Type)
	f.Set("Legal", o.Legal)

	f.Set("Area", fmt.Sprint(o.Area))
	f.Set("Bedroom", fmt.Sprint(o.Bedroom))
	f.Set("Toilet", fmt.Sprint(o.Toilet))
	f.Set("Floor", fmt.Sprint(o.Floor))
	f.Set("FrontWidth", fmt.Sprint(o.FrontWidth))
	f.Set("StreetWidth", fmt.Sprint(o.StreetWidth))
	f.Set("PavementWidth", fmt.Sprint(o.PavementWidth))

	return f
}

func (o *Attachment) Form() *form.Form {
	f := form.New(nil)
	f.Set("Title", o.Title)
	f.Set("ContentType", o.ContentType)
	f.Set("MineType", o.MineType)
	f.Set("Link", o.Link)
	f.Set("VideoThumbnail", o.VideoThumbnail)

	f.Set("ProductID", fmt.Sprint(o.Product.ID))
	f.Set("Width", fmt.Sprint(o.Width))
	f.Set("Height", fmt.Sprint(o.Height))
	f.Set("Size", fmt.Sprint(o.Size))
	f.Set("VideoLength", fmt.Sprint(o.VideoLength))

	return f
}

type Log struct{}
