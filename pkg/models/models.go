package models

import "time"

type User struct {
	ID        int
	Email     string
	Phone     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
