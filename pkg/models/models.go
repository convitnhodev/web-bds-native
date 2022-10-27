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
