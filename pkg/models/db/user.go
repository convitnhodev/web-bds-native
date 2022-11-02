package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/helper"
	"github.com/deeincom/deeincom/pkg/models"
)

type UserModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var userColumes = []string{
	"users.id",
	"users.first_name",
	"users.last_name",
	"users.email",
	"users.phone",
	"users.password",
	"users.roles",
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
		pq.Array(&o.Roles),
	); err != nil {
		return errors.Wrap(err, "scanUser")
	}

	return nil
}

func (m *UserModel) Auth(f *form.Form) (*models.User, error) {
	q := `
		select id, password from users where phone = $1
	`
	var id int
	var hashed []byte
	row := m.DB.QueryRow(q, f.Get("Phone"))
	if err := row.Scan(&id, &hashed); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(hashed, []byte(f.Get("Password"))); err != nil {
		return nil, err
	}

	return &models.User{ID: id}, nil
}

func (m *UserModel) Create(f *form.Form) (*models.User, error) {
	q := `
	insert into users (
		first_name,
		last_name,
		email,
		phone,
		password,
		email_token,
		phone_token
	) values (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	) returning id
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
		helper.RandString(6),
		helper.RandString(6),
	)
	o := new(models.User)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	return o, nil
}

func (m *UserModel) ID(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("err_id_empty")
	}
	q := m.query(`where users.id = $1`)
	row := m.DB.QueryRow(q, id)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.ID")
	}
	return o, nil
}

func (m *UserModel) GetByEmailToken(token string) (*models.User, error) {
	q := m.query(`where users.email_token = $1`)
	row := m.DB.QueryRow(q, token)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.GetByEmailToken")
	}
	return o, nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	q := m.query(`where users.email = $1`)
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
	q := m.query("order by id desc")
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

func (m *UserModel) LogSendVerifyEmail(user *models.User) error {
	q := `
		update users
			set email = $2,
			send_verified_email_at = now()
		where id = $1
	`
	_, err := m.DB.Exec(q, user.ID, user.Email)
	return err
}

func (m *UserModel) LogSendVerifyPhone(user *models.User) error {
	q := `
		update users
			set phone = $2,
			send_verified_phone_at = now()
		where id = $1
	`
	_, err := m.DB.Exec(q, user.ID, user.Phone)
	return err
}
