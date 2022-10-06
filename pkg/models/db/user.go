package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/deeincom/deeincom/pkg/models"
)

type UserModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

type scanner interface {
	Scan(dest ...interface{}) error
}

var cols = []string{"user_id", "displayname"}

func (m *UserModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM users`, strings.Join(cols, ","))
}

func scanUser(r scanner, user *models.User) error {
	if err := r.Scan(
		&user.ID,
		&user.Name,
	); err != nil {
		return errors.Wrap(err, "scanUser")
	}

	return nil
}

// Create a new user
func (m *UserModel) Create(u *models.User) (*models.User, error) {
	q := `
	insert into users (
		displayname,
		password
	) values (
		$1,
		$2
	) returning user_id
	`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	row := m.DB.QueryRow(q,
		u.Name,
		string(hashedPassword),
	)

	if err := row.Scan(&u.ID); err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	return u, nil
}

// ID return user by ID
func (m *UserModel) ID(id int) (*models.User, error) {
	q := m.query(`where users.user_id = $1`)
	row := m.DB.QueryRow(q, id)
	u := &models.User{}
	if err := scanUser(row, u); err != nil {
		return nil, errors.Wrap(err, "user.ID")
	}
	return u, nil
}
