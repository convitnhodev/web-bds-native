package repositories

import (
	"github.com/deeincom/deeincom/app/models"
	"github.com/upper/db/v4"
	"time"
)

type UserRepository repository

func (u *UserRepository) CreateAnUser(user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := u.r.db.Collection(models.UserTable).InsertReturning(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByPhoneNumber(phone string) (*models.User, error) {
	var user models.User
	if err := u.Find(db.Cond{"phone_number": phone}).One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
