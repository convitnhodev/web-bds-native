package models

import (
	"time"
)

var UserVerifyTable = "user_verifies"

type UserVerify struct {
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Token       string    `json:"token" db:"token"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
