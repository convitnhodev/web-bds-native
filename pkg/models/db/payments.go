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

type PaymentModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var paymentColumes = []string{
	"payments.id",
	"payments.invoice_id",
	"payments.amount",
	"payments.method",
	"payments.pay_type",
	"payments.tx_type",
	"payments.status",
	"payments.created_at",
	"payments.updated_at",
}

func scanPayment(r scanner, o *models.Payment) error {
	if err := r.Scan(
		&o.ID,
		&o.InvoiceId,
		&o.Amount,
		&o.Method,
		&o.PayType,
		&o.TxType,
		&o.Status,
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

func (m *PaymentModel) ID(id string) (*models.Payment, error) {
	q := m.query(`where payments.id = $1`)
	row := m.DB.QueryRow(q, id)

	o := new(models.Payment)
	if err := scanPayment(row, o); err != nil {
		return nil, errors.Wrap(err, "PaymentModel.ID")
	}

	return o, nil
}

func (m *PaymentModel) InvoiceID(invoiceId int) ([]*models.Payment, error) {
	q := m.query("WHERE payments.invoice_id = $1")

	rows, err := m.DB.Query(q, invoiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Payment{}
	for rows.Next() {
		o := &models.Payment{}
		if err := scanPayment(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *PaymentModel) Checkout(
	tx *sql.Tx,
	ctx context.Context,
	invoiceId int,
	amount int64,
	method string,
	payType string,
) (*models.Payment, error) {
	q := `
		INSERT INTO payments (invoice_id, amount, pay_type)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	row := tx.QueryRowContext(
		ctx,
		q,
		invoiceId,
		amount,
		payType,
	)

	o := new(models.Payment)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "PaymentModel.Checkout")
	}

	return o, nil
}

func (m *PaymentModel) UpdatePostData(
	tx *sql.Tx,
	ctx context.Context,
	paymentId int,
	postData string,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			post_data = $2
		WHERE id = $1
	`

	_, err := tx.ExecContext(
		ctx,
		q,
		paymentId,
		postData,
	)

	return err
}

func (m *PaymentModel) UpdatePaymentCallback(
	tx *sql.Tx,
	ctx context.Context,
	paymentId int,
	isSuccess bool,
	paymentData string,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			status = $2,
			recipition_data = $3
		WHERE id = $1
	`

	status := "failed"
	if isSuccess {
		status = "success"
	}

	_, err := tx.ExecContext(
		ctx,
		q,
		paymentId,
		status,
		paymentData,
	)

	return err
}
