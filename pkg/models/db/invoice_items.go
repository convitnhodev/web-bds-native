package db

import (
	"database/sql"
	"fmt"
	"log"
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
	"products.id",
	"products.title",
	"products.slug",
}

func scanInvoiceItem(r scanner, o *models.InvoiceItem) error {
	o.Product = models.Product{}
	if err := r.Scan(
		&o.ID,
		&o.InvoiceId,
		&o.ProductId,
		&o.Quatity,
		&o.CostPerSlot,
		&o.Amount,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Product.ID,
		&o.Product.Title,
		&o.Product.Slug,
	); err != nil {
		return errors.Wrap(err, "scanInvoice")
	}

	return nil
}

func (m *InvoiceItemModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM invoice_items JOIN products ON invoice_items.product_id = products.id %s`, strings.Join(invoiceItemsColumes, ","), s)
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
	q := m.query("WHERE invoice_items.invoice_id = $1")

	rows, err := m.DB.Query(q, invoiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.InvoiceItem{}
	for rows.Next() {
		o := &models.InvoiceItem{}
		if err := scanInvoiceItem(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *InvoiceItemModel) Buy(invoiceId int, productId int, quatity int, cost int, amount int) (*models.InvoiceItem, error) {
	q := `
	INSERT INTO public.invoice_items
	(invoice_id, product_id, quatity, cost_per_slot, amount)
	VALUES($1, $2, $3, $4, $5);
	RETURNING id`

	row := m.DB.QueryRow(q,
		invoiceId,
		productId,
		quatity,
		cost,
		amount,
	)

	o := new(models.InvoiceItem)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "InvoiceItemModel.Buy")
	}

	return o, nil
}
