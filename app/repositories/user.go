package repositories

import "github.com/upper/db/v4"

const (
	userTable = "users"
)

type User repository

func (u *User) IsPhoneExist(pn string) (bool, error) {
	return u.r.db.Collection(userTable).Find(db.Cond{"phone_number": db.Is(pn)}).Exists()
}
