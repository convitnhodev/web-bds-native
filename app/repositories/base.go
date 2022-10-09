package repositories

import "github.com/upper/db/v4"

type Repository struct {
	common repository
	db     db.Session

	User *User
}

type repository struct {
	r *Repository
}

func New(db db.Session) *Repository {
	r := &Repository{
		db: db,
	}
	r.User = (*User)(&r.common)
	return r
}
