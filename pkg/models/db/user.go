package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
)

type UserModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var userColumes = []string{
	"users.user_id",
	"users.user_first_name",
	"users.user_last_name",
	"users.user_email",
	"users.user_phone",
	"users.user_password",
}

func (m *UserModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM users`, strings.Join(userColumes, ","))
}

func scanUser(r scanner, o *models.User) error {
	if err := r.Scan(
		&o.ID,
		&o.FirstName,
		&o.LastName,
		&o.Email,
		&o.Phone,
		&o.Password,
	); err != nil {
		return errors.Wrap(err, "scanUser")
	}

	return nil
}

func (m *UserModel) Create(f *form.Form) (*models.User, error) {
	q := `
	insert into users (
		user_first_name,
		user_last_name,
		user_email,
		user_phone,
		user_password
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
	) returning user_id
	`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Get("Password")), 12)
	if err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	row := m.DB.QueryRow(q,
		f.Get("FirstName"),
		f.Get("LastName"),
		f.Get("Email"),
		f.Get("Phone"),
		string(hashedPassword),
	)
	o := new(models.User)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	return o, nil
}

func (m *UserModel) ID(id int) (*models.User, error) {
	q := m.query(`where users.user_id = $1`)
	row := m.DB.QueryRow(q, id)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.ID")
	}
	return o, nil
}

func (m *UserModel) GetByEmailToken(token string) (*models.User, error) {
	return nil, nil
}

func (m *UserModel) AddRole(id int, role string) error {
	return nil
}
