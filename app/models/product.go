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

func (p *Product) GetCategoryDisplay() string {
	categories := map[string]string{
		"1":  "Căn hộ chung cư",
		"2":  "Nhà riêng",
		"3":  "Nhà biệt thự, liền kề",
		"4":  "Nhà mặt phố",
		"5":  "Shop house, nhà phố thương mại",
		"6":  "Đất nền dự án",
		"7":  "Đất",
		"8":  "Trang trại, khu nghỉ dưỡng",
		"9":  "Condotel",
		"10": "Kho, nhà xưởng",
		"99": "Các loại hình bất động sản khác",
	}

	if v, ok := categories[p.CategoryID]; ok {
		return v
	}
	return categories["99"]
}

func (p *Product) GetDocumentTypeDisplay() string {
	docTypes := map[string]string{
		"1": "Sổ đỏ, sổ hồng",
		"2": "Hợp đồng mua bán",
		"3": "Đang chờ sổ",
		"4": "Khác",
	}

	if v, ok := docTypes[p.DocumentType]; ok {
		return v
	}
	return docTypes["4"]
}

func (p *Product) GetHouseDirectionDisplay() string {
	docTypes := map[string]string{
		"1": "Bắc",
		"2": "Nam",
		"3": "Đông",
		"4": "Tây",
		"5": "Đông Bắc",
		"6": "Đông Nam",
		"7": "Tây Bắc",
		"8": "Tây Nam",
	}

	if v, ok := docTypes[p.HouseDirection]; ok {
		return v
	}
	return "unknown"
}

func (p *Product) GetBalconyDirectionDisplay() string {
	docTypes := map[string]string{
		"1": "Bắc",
		"2": "Nam",
		"3": "Đông",
		"4": "Tây",
		"5": "Đông Bắc",
		"6": "Đông Nam",
		"7": "Tây Bắc",
		"8": "Tây Nam",
	}

	if v, ok := docTypes[p.BalconyDirection]; ok {
		return v
	}
	return "unknown"
}
