package repositories

import (
	"fmt"
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

func (u *UserRepository) UpdateVerified(phone string) error {
	return u.r.db.Collection(models.UserTable).Find(db.Cond{"phone_number": phone}).Update(map[string]time.Time{
		"verified_at": time.Now(),
		"updated_at":  time.Now(),
	})
}

func (u *UserRepository) CreateVerifyCode(phone string, token string) error {
	vc := &models.UserVerify{
		PhoneNumber: phone,
		Token:       token,
		CreatedAt:   time.Now(),
	}
	if _, err := u.r.db.Collection(models.UserVerifyTable).Insert(vc); err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetVerifyCode(phone string, token string) (*models.UserVerify, error) {
	var uv models.UserVerify
	if err := u.r.db.Collection(models.UserVerifyTable).Find(db.And(
		db.Cond{"phone_number": phone},
		db.Cond{"token": token},
		db.Cond{"created_at >": time.Now().Add(-2 * time.Minute)},
	)).One(&uv); err != nil {
		return nil, err
	}
	fmt.Printf("row: %+v\n", uv)
	return &uv, nil
}
