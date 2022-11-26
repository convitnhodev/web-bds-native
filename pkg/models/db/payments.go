package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type PaymentModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var paymentColumes = []string{
	"payments.id",
	"payments.invoice_id",
	"payments.amount",
	"payments.menthod",
	"payments.pay_type",
	"payments.tx_type",
	"payments.created_at",
	"payments.updated_at",
}

func scanPayment(r scanner, o *models.Payment) error {
	if err := r.Scan(
		&o.ID,
		&o.InvoiceId,
		&o.Amount,
		&o.Menthod,
		&o.PayType,
		&o.TxType,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanPayment")
	}

	return nil
}

func (m *PaymentModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM payments %s`, strings.Join(paymentColumes, ","), s)
}

func (m *PaymentModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM payments %s`, s)
}

func (m *PaymentModel) InvoiceID(invoiceId int) ([]*models.Payment, error) {
	return nil, nil
}
