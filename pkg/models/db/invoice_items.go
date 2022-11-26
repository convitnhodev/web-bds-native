package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type InvoiceItemModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var invoiceItemsColumes = []string{
	"invoice_items.id",
	"invoice_items.invoice_id",
	"invoice_items.product_id",
	"invoice_items.quatity",
	"invoice_items.cost_per_slot",
	"invoice_items.amount",
	"invoice_items.created_at",
	"invoice_items.updated_at",
}

func scanInvoiceItems(r scanner, o *models.InvoiceItem) error {
	if err := r.Scan(
		&o.ID,
		&o.InvoiceId,
		&o.ProductId,
		&o.Quatity,
		&o.CostPerSlot,
		&o.Amount,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanInvoice")
	}

	return nil
}

func (m *InvoiceItemModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM invoice_items %s`, strings.Join(invoiceItemsColumes, ","), s)
}

func (m *InvoiceItemModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM invoice_items %s`, s)
}

func (m *InvoiceItemModel) CountByProduct(productId int) (*int, error) {
	q := m.count("where product_id = $1")

	count := 0
	row := m.DB.QueryRow(q, productId)
	if err := row.Scan(&count); err != nil {
		return nil, err
	}

	return &count, nil
}

func (m *InvoiceItemModel) InvoiceID(invoiceId int) ([]*models.InvoiceItem, error) {
	return nil, nil
}
