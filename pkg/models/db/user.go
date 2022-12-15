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
	"users.created_at",
	"users.updated_at",
	"users.send_verified_email_at",
	"users.send_verified_phone_at",
	"users.reset_pwd_token",
	"users.rpt_expired_at",
	"users.partner_status",
	"users.last_kyc_status",
}

func (m *UserModel) HashPassword(s string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), 12)
	return string(hashedPassword), err
}

func (m *UserModel) CompareHashAndPassword(hashed string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err
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
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.SendVerifiedEmailAt,
		&o.SendVerifiedPhoneAt,
		&o.ResetPasswordToken,
		&o.RPTExpiredAt,
		&o.PartnerStatus,
		&o.LastKYCStatus,
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
	var hashed string
	row := m.DB.QueryRow(q, f.Get("Phone"))
	if err := row.Scan(&id, &hashed); err != nil {
		return nil, err
	}

	if err := m.CompareHashAndPassword(hashed, f.Get("Password")); err != nil {
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
	hashedPassword, err := m.HashPassword(f.Get("Password"))

	if err != nil {
		return nil, errors.Wrap(err, "UserModel.Create")
	}

	row := m.DB.QueryRow(q,
		f.Get("FirstName"),
		f.Get("LastName"),
		f.Get("Email"),
		f.Get("Phone"),
		hashedPassword,
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

func (m *UserModel) GetByPhoneToken(token string) (*models.User, error) {
	q := m.query(`where users.phone_token = $1`)
	row := m.DB.QueryRow(q, token)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.GetByPhoneToken")
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

func (m *UserModel) GetByPhone(phone string) (*models.User, error) {
	q := m.query(`where users.phone = $1`)
	row := m.DB.QueryRow(q, phone)
	o := new(models.User)
	if err := scanUser(row, o); err != nil {
		return nil, errors.Wrap(err, "user.GetByPhone")
	}
	return o, nil
}

func (m *UserModel) AddRole(user *models.User, role string) error {
	for _, s := range user.Roles {
		if s == role {
			return errors.New("err_duplicated_role")
		}
	}
	q := `update users set roles = array_append(roles, $2) where id = $1`
	_, err := m.DB.Exec(q, user.ID, role)
	return err
}

func (m *UserModel) Find(kycStatus string, partnerStatus string) ([]*models.User, error) {
	q := m.query("order by id desc")

	if kycStatus != "" {
		q = m.query(fmt.Sprintf("WHERE last_kyc_status = '%s' order by id desc", kycStatus))
	}

	if partnerStatus != "" {
		q = m.query(fmt.Sprintf("WHERE partner_status = '%s' order by id desc", partnerStatus))
	}

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
	user.EmailToken = helper.RandString(6)
	q := `
		update users set
			email = $2,
			send_verified_email_at = now(),
			email_token = $3
		where id = $1
	`
	_, err := m.DB.Exec(q, user.ID, user.Email, user.EmailToken)
	return err
}

func (m *UserModel) LogSendVerifyPhone(user *models.User) error {
	user.PhoneToken = helper.RandString(6)

	q := `
		update users set
			phone = $2,
			send_verified_phone_at = now(),
			phone_token = $3
		where id = $1
	`
	_, err := m.DB.Exec(q, user.ID, user.Phone, user.PhoneToken)
	return err
}

func (m *UserModel) UpdateKYCStatus(userId string, status string) error {
	q := `
		UPDATE users SET
			updated_at = now(),
			last_kyc_status = $2,
			roles = ARRAY(SELECT DISTINCT e FROM UNNEST(ARRAY_APPEND(roles, $3)) AS t(e) WHERE e != '')
		WHERE id = $1
	`

	newRole := ""
	if status == "approved_kyc" {
		newRole = "verified_id"
	}

	_, err := m.DB.Exec(q,
		userId,
		status,
		newRole,
	)

	return err
}

func (m *UserModel) UpdatePartnerStatus(userId string, status string) error {
	q := `
		UPDATE users SET
			updated_at = now(),
			partner_status = $2,
			roles = ARRAY(SELECT DISTINCT e FROM UNNEST(ARRAY_APPEND(roles, $3)) AS t(e) WHERE e != '')
		WHERE id = $1
	`

	newRole := ""
	if status == "approved" {
		newRole = "deein_partner"
	}

	_, err := m.DB.Exec(q,
		userId,
		status,
		newRole,
	)

	return err
}

func (m *UserModel) ResetPasswordByPhone(phone string, token string) error {
	hashedToken, err := m.HashPassword(token)
	if err != nil {
		return err
	}

	q := `
		UPDATE users SET
			updated_at = now(),
			reset_pwd_token = $2,
			rpt_expired_at = now() + (15 * interval '1 minute')
		WHERE phone = $1
	`
	_, qerr := m.DB.Exec(q,
		phone,
		hashedToken,
	)

	return qerr
}

func (m *UserModel) UpdateNewPassword(userId string, password string) error {
	hashedPassword, err := m.HashPassword(password)

	if err != nil {
		return err
	}

	q := `
		UPDATE users SET
			updated_at = now(),
			password = $2,
			rpt_expired_at = null,
			reset_pwd_token = ''
		WHERE id = $1
	`

	_, qerr := m.DB.Exec(q,
		userId,
		hashedPassword,
	)

	return qerr
}
