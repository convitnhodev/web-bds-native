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

type PartnerModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var partnerColumes = []string{
	"partner.id",
	"partner.user_id",
	"partner.message",
	"partner.cv_link",
	"partner.status",
	"partner.feedback",
	"partner.approved_by",
	"partner.rejected_by",
	"partner.created_at",
	"partner.updated_at",
}

func scanPartner(r scanner, o *models.Partner) error {
	if err := r.Scan(
		&o.ID,
		&o.UserId,
		&o.Message,
		&o.CVLink,
		&o.Status,
		&o.Feedback,
		&o.ApprovedBy,
		&o.RejectedBy,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanFile")
	}

	return nil
}

func (m *PartnerModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM partner %s`, strings.Join(partnerColumes, ","), s)
}

func (m *PartnerModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM partner %s`, s)
}

func (m *PartnerModel) User(userId string) ([]*models.Partner, error) {
	q := m.query(fmt.Sprintf("WHERE partner.user_id = '%s' ORDER BY partner.id DESC LIMIT 1", userId))

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Partner{}
	for rows.Next() {
		o := &models.Partner{}
		if err := scanPartner(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}

	return list, nil
}

func (m *PartnerModel) SubmitPartner(userId string, f *form.Form) error {
	q := `
		INSERT INTO partner (user_id, message, cv_link)
		VALUES ($1, $2, $3)
	`
	_, err := m.DB.Exec(q,
		userId,
		f.Get("Message"),
		f.Get("CVLink"),
	)

	return err
}

func (m *PartnerModel) FeedbackPartner(partnerId string, userId string, status string, reason string) error {
	q := `
		UPDATE partner SET
			updated_at = now(),
			approved_by = $2,
			status = $3,
			feedback = $4
		WHERE id = $1
	`

	if status == "rejected" {
		q = `
			UPDATE partner SET
				updated_at = now(),
				rejected_by = $2,
				status = $3,
				feedback = $4
			WHERE id = $1
		`
	}

	_, err := m.DB.Exec(q,
		partnerId,
		userId,
		status,
		reason,
	)

	return err
}
