package repositories

import (
	"github.com/deeincom/deeincom/app/models"
	"github.com/upper/db/v4"
)

type repository struct {
	r *Repository
	db.Collection
}

type Repository struct {
	db     db.Session
	common repository

	User *UserRepository
}

func New(db db.Session) *Repository {
	r := &Repository{
		db: db,
	}
	r.common.r = r
	r.User = (*UserRepository)(&repository{r: r, Collection: r.db.Collection(models.UserTable)})
	return r
}
