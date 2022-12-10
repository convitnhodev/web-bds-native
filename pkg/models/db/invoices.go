package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type InvoiceModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var invoiceColumes = []string{
	"invoices.id",
	"invoices.user_id",
	"invoices.status",
	"invoices.invoice_synced_at",
	"invoices.invoice_serect",
	"invoices.created_at",
	"invoices.updated_at",
	"users.id",
	"users.first_name",
	"users.last_name",
	"users.email",
	"users.phone",
}

func scanInvoice(r scanner, o *models.Invoice, includeUser bool) error {
	if includeUser {
		o.User = models.User{}
		if err := r.Scan(
			&o.ID,
			&o.UserId,
			&o.Status,
			&o.InvoiceSyncedAt,
			&o.InvoiceSerect,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.User.ID,
			&o.User.FirstName,
			&o.User.LastName,
			&o.User.Email,
			&o.User.Phone,
		); err != nil {
			return errors.Wrap(err, "scanInvoice")
		}
	} else {
		if err := r.Scan(
			&o.ID,
			&o.UserId,
			&o.Status,
			&o.InvoiceSyncedAt,
			&o.InvoiceSerect,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return errors.Wrap(err, "scanInvoice")
		}
	}

	return nil
}

func (m *InvoiceModel) query(s string, includeUser bool) string {
	if includeUser {
		return fmt.Sprintf(`SELECT %s FROM invoices JOIN users ON invoices.user_id = users.id %s`, strings.Join(invoiceColumes, ","), s)
	} else {
		return fmt.Sprintf(`SELECT %s FROM invoices %s`, strings.Join(invoiceColumes, ","), s)
	}
}

func (m *InvoiceModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM invoices %s`, s)
}

func (m *InvoiceModel) ID(id string) (*models.Invoice, error) {
	q := m.query("WHERE invoices.id = $1", true)

	row := m.DB.QueryRow(q, id)
	o := new(models.Invoice)
	if err := scanInvoice(row, o, true); err != nil {
		return nil, errors.Wrap(err, "InvoiceModel.ID")
	}

	return o, nil
}

func (m *InvoiceModel) Buy(tx *sql.Tx, ctx context.Context, userId int, serect string) (*models.Invoice, error) {
	q := `
	INSERT INTO invoices (user_id, status, invoice_serect)
	VALUES($1, 'open'::invoice_status, $2)
	RETURNING id`

	row := tx.QueryRowContext(
		ctx,
		q,
		userId,
		serect,
	)

	o := new(models.Invoice)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "Invoice.Buy")
	}

	return o, nil
}

func (m *InvoiceModel) Find(productId int) ([]*models.Invoice, error) {
	q := fmt.Sprintf(`
		SELECT
			%s
		FROM invoices
		WHERE invoices.id IN (
			SELECT
				DISTINCT invoice_id
			FROM invoice_items
			WHERE invoice_items.product_id = $1
		)
	`, strings.Join(invoiceColumes[0:7], ","))

	count := `
		SELECT
			count(*)
		FROM invoices
		WHERE invoices.id IN (
			SELECT
				DISTINCT invoice_id
			FROM invoice_items
			WHERE invoice_items.product_id = $1
		)
	`

	if err := m.Pagination.Count(count, productId); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(m.Pagination.Generate(q), productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Invoice{}
	for rows.Next() {
		o := &models.Invoice{}
		if err := scanInvoice(rows, o, false); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *InvoiceModel) UpdatePaymentCallback(
	tx *sql.Tx,
	ctx context.Context,
	invoiceId int,
	transactionTs int,
) error {
	q := `
		UPDATE invoices
		SET updated_at = now(),
			invoice_synced_at = to_timestamp($2),
			status = 'deposit'
		WHERE id = $1
	`

	_, err := tx.ExecContext(
		ctx,
		q,
		invoiceId,
		transactionTs,
	)

	return err
}
