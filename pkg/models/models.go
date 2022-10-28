package models

import "time"

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
	Title  string
	Type   string
	URL    string
	Width  int
	Height int
	Size   int
	Length int
}

type Product struct {
	ID                   int
	Title                string
	Short                string
	Full                 string
	City                 string
	District             string
	Ward                 string
	AddressNumber        string
	Street               string
	HouseDirection       string
	BalconyDirection     string
	BusinessAdvantage    string
	FinancialPlan        string
	Furniture            string
	GeneralPurchaseTerms string
	Area                 int
	Bedroom              int
	Toilet               int
	Floor                int
	FrontWidth           int
	StreetWidth          int
	PavementWidth        int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
