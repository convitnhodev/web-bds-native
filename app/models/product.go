package models

import (
	"github.com/upper/db/v4"
	"time"
)

var ProductTable = "products"

type Product struct {
	ID                int64     `json:"id" db:"id"`
	Title             string    `json:"title" db:"title"`
	CategoryID        string    `json:"category_id" db:"category_id"`
	Description       string    `json:"description" db:"description"`
	City              string    `json:"city" db:"city"`
	District          string    `json:"district" db:"district"`
	Ward              string    `json:"ward" db:"ward"`
	AddressNumber     string    `json:"address_number" db:"address_number"`
	Street            string    `json:"street" db:"street"`
	Area              int64     `json:"area" db:"area"`
	DocumentType      string    `json:"document_type" db:"document_type"`
	DocumentProof     string    `json:"document_proof" db:"document_proof"`
	Bedroom           int64     `json:"bedroom" db:"bedroom"`
	Toilet            int64     `json:"toilet" db:"toilet"`
	Floor             int64     `json:"floor" db:"floor"`
	HouseDirection    string    `json:"house_direction" db:"house_direction"`
	BalconyDirection  string    `json:"balcony_direction" db:"balcony_direction"`
	FrontWidth        int64     `json:"front_width" db:"front_width"`
	StreetWidth       int64     `json:"street_width" db:"street_width"`
	PavementWidth     int64     `json:"pavement_width" db:"pavement_width"`
	BusinessAdvantage string    `json:"business_advantage" db:"business_advantage"`
	FinancialPlan     string    `json:"financial_plan" db:"financial_plan"`
	Furniture         string    `json:"furniture" db:"furniture"`
	IsActivated       bool      `json:"is_activated" db:"is_activated"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Product) Store(sess db.Session) db.Store {
	return sess.Collection(ProductTable)
}

var _ = db.Record(&Product{})
