package db

import (
	"database/sql"
	"fmt"
	"log"
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
	return fmt.Sprintf(`SELECT %s FROM users %s`, strings.Join(userColumes, ","), s)
}

func (m *UserModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM users %s`, s)
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
	q := m.query(`where users.user_email_token = $1`)
	row := m.DB.QueryRow(q, token)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.GetByEmailToken")
	}
	return o, nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	q := m.query(`where users.user_email = $1`)
	row := m.DB.QueryRow(q, email)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.GetByEmail")
	}
	return o, nil
}

func (m *UserModel) AddRole(id int, role string) error {
	return nil
}

func (m *UserModel) Find() ([]*models.User, error) {
	q := m.query("order by user_id desc")
	count := m.count("")

	if err := m.Pagination.Count(count); err != nil {
		return nil, err
	}
	rows, err := m.DB.Query(m.Pagination.Generate(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.User{}
	for rows.Next() {
		o := &models.User{}
		if err := scanUser(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}
