package models

import (
	"github.com/upper/db/v4"
	"time"
)

var UserTable = "users"

type User struct {
	ID          int64     `json:"id" db:"id,omitempty"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       *string   `json:"email,omitempty" db:"email,omitempty"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	VerifiedAt  time.Time `json:"verified_at,omitempty" db:"email_verified_at,omitempty"`
	Password    string    `json:"password" db:"password"`
	IsActivated bool      `json:"is_activated" db:"is_activated"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (book *User) Store(sess db.Session) db.Store {
	return sess.Collection(UserTable)
}

var _ = db.Record(&User{})
