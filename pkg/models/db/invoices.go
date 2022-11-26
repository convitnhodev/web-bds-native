package db

import (
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
}

func scanInvoice(r scanner, o *models.Invoice) error {
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

	return nil
}

func (m *InvoiceModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM invoices %s`, strings.Join(invoiceColumes, ","), s)
}

func (m *InvoiceModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM invoices %s`, s)
}

func (m *InvoiceModel) ID(id string) (*models.Invoice, error) {
	q := m.query("WHERE invoices.id = $1")

	row := m.DB.QueryRow(q, id)
	o := new(models.Invoice)
	if err := scanInvoice(row, o); err != nil {
		return nil, errors.Wrap(err, "InvoiceModel.ID")
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
	`, strings.Join(invoiceColumes, ","))

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
		if err := scanInvoice(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}
