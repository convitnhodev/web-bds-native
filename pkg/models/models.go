package models

import "time"

type User struct {
	ID        int
	Email     string
	Role      string
	Password  string `json:"password"`
	Name      string
	CreatedAt time.Time
}
