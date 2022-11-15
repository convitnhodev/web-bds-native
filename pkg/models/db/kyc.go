package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type KYCModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var kycColumes = []string{
	"kyc.id",
	"kyc.user_id",
	"kyc.front_identity_card",
	"kyc.back_identity_card",
	"kyc.selfie_image",
	"kyc.feedback",
	"kyc.approved_by",
	"kyc.rejected_by",
	"kyc.status",
	"kyc.created_at",
	"kyc.updated_at",
}

func scanKYC(r scanner, o *models.KYC) error {
	if err := r.Scan(
		&o.ID,
		&o.UserId,
		&o.FrontIdentityCard,
		&o.BackIdentityCard,
		&o.SelfieImage,
		&o.Feedback,
		&o.ApprovedBy,
		&o.RejectedBy,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanFile")
	}

	return nil
}

func (m *KYCModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM kyc %s`, strings.Join(kycColumes, ","), s)
}

func (m *KYCModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM files %s`, s)
}

func (m *KYCModel) User(userId string) ([]*models.KYC, error) {
	q := m.query(fmt.Sprintf("WHERE kyc.user_id = '%s' ORDER BY kyc.id DESC LIMIT 1", userId))

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.KYC{}
	for rows.Next() {
		o := &models.KYC{}
		if err := scanKYC(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}

	return list, nil
}

func (m *KYCModel) SubmitKYC(userId string, f *form.Form) error {
	q := `
		INSERT INTO kyc (userId, front_identity_card, back_identity_card, selfie_image)
		VALUES ($1, $2, $3, $4)
	`
	_, err := m.DB.Exec(q,
		userId,
		f.Get("FrontIdentityCard"),
		f.Get("BackIdentityCard"),
		f.Get("SelfieImage"),
	)

	return err
}

func (m *KYCModel) FeedbackKYC(kycId string, userId string, status string, reason string) error {
	q := `
		UPDATE kyc SET
			updated_at = now(),
			approved_by = $2,
			status = $3,
			feedback = $4
		WHERE id = $1
	`

	if status == "rejected_kyc" {
		q = `
			UPDATE kyc SET
				updated_at = now(),
				rejected_by = $2,
				status = $3,
				feedback = $4
			WHERE id = $1
		`
	}

	_, err := m.DB.Exec(q,
		kycId,
		userId,
		status,
		reason,
	)

	return err
}
