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
	"invoices.total_amount",
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
			&o.TotalAmount,
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
			&o.TotalAmount,
		); err != nil {
			return errors.Wrap(err, "scanInvoice")
		}
	}

	return nil
}

func (m *InvoiceModel) query(s string, includeUser bool) string {
	userJoinSql := ""
	columes := invoiceColumes[0:8]
	if includeUser {
		columes = invoiceColumes
		userJoinSql = "LEFT JOIN users ON invoices.user_id = users.id"
	}

	return fmt.Sprintf(`SELECT %s FROM invoices %s %s`, strings.Join(columes, ","), userJoinSql, s)
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

func (m *InvoiceModel) Buy(tx *sql.Tx, ctx context.Context, userId int, serect string, totalAmount int) (*models.Invoice, error) {
	q := `
	INSERT INTO invoices (user_id, status, invoice_serect, total_amount)
	VALUES($1, 'open'::invoice_status, $2, $3)
	RETURNING id`

	row := tx.QueryRowContext(
		ctx,
		q,
		userId,
		serect,
		totalAmount,
	)

	o := new(models.Invoice)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "Invoice.Buy")
	}

	return o, nil
}

func (m *InvoiceModel) Find(productId int) ([]*models.Invoice, error) {
	q := m.query(`
		WHERE invoices.id IN (
			SELECT
				DISTINCT invoice_id
			FROM invoice_items
			WHERE invoice_items.product_id = $1
		)
	`, false)

	count := m.count(`
		WHERE invoices.id IN (
			SELECT
				DISTINCT invoice_id
			FROM invoice_items
			WHERE invoice_items.product_id = $1
		)
	`)

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
	isSuccess bool,
) error {
	q := `
		UPDATE invoices
		SET updated_at = now(),
			invoice_synced_at = to_timestamp($2),
			status = $3`

	status := "slot_canceled"
	if isSuccess {
		status = "deposit"
		q += `
		WHERE id = $1`
	} else {
		q += `,
			slot_canceled_by = user_id
		WHERE id = $1
		`
	}

	_, err := tx.ExecContext(
		ctx,
		q,
		invoiceId,
		transactionTs,
		status,
	)

	return err
}

func (m *InvoiceModel) Refund(
	id int,
	userId int,
) error {
	q := `
	UPDATE invoices
	SET updated_at = now(),
		status = $2,
		slot_canceled_by = $3
	WHERE id = $1`

	_, err := m.DB.Exec(q,
		id,
		"refund",
		userId,
	)

	return err
}

func (m *InvoiceModel) UpdateStatus(
	id int,
	status string,
) error {
	q := `
	UPDATE invoices
	SET updated_at = now(),
		status = $2
	WHERE id = $1`

	_, err := m.DB.Exec(q,
		id,
		status,
	)

	return err
}

func (m *InvoiceModel) CollectMoney(
	tx *sql.Tx,
	ctx context.Context,
	invoiceId int,
) error {
	q := `
		UPDATE invoices
		SET updated_at = now(),
			status = $2
		WHERE id = $1`

	_, err := tx.ExecContext(
		ctx,
		q,
		invoiceId,
		"collecting",
	)

	return err
}
