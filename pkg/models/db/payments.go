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
	"payments.appotapay_trans_id",
	"payments.refund_id",
	"payments.refund_response",
	"payments.transaction_at",
	"payments.refund_at",
	"payments.appotapay_account_no",
	"payments.appotapay_account_name",
	"payments.appotapay_bank_code",
	"payments.appotapay_bank_name",
	"payments.appotapay_bank_branch",
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
		&o.AppotapayTransId,
		&o.RefundId,
		&o.RefundResponse,
		&o.TransactionAt,
		&o.RefundAt,
		&o.AppotapayAccountNo,
		&o.AppotapayAccountName,
		&o.AppotapayBankCode,
		&o.AppotapayBankName,
		&o.AppotapayBankBranch,
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
	amount int,
	method string,
	payType string,
) (*models.Payment, error) {
	q := `
		INSERT INTO payments (invoice_id, amount, method, pay_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	row := tx.QueryRowContext(
		ctx,
		q,
		invoiceId,
		amount,
		method,
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
	resData string,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			post_data = $2,
			response_data = $3
		WHERE id = $1
	`

	_, err := tx.ExecContext(
		ctx,
		q,
		paymentId,
		postData,
		resData,
	)

	return err
}

func (m *PaymentModel) UpdatePaymentCallback(
	tx *sql.Tx,
	ctx context.Context,
	paymentId int,
	isSuccess bool,
	paymentData string,
	apptotapayTransId string,
	transactionAt int,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			status = $2,
			recipition_data = $3,
			appotapay_trans_id = $4,
			transaction_at = to_timestamp($5)
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
		apptotapayTransId,
		transactionAt,
	)

	return err
}

func (m *PaymentModel) Refund(
	id int,
	refundResponse string,
	refundId string,
	refundAt int,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			status = $2,
			refund_response = $3,
			refund_id = $4,
			refund_at = to_timestamp($5)
		WHERE id = $1`

	_, err := m.DB.Exec(q,
		id,
		"refund",
		refundResponse,
		refundId,
		refundAt,
	)

	return err
}

func (m *PaymentModel) UpdateBankAcount(
	tx *sql.Tx,
	ctx context.Context,
	paymentId int,
	billCode string,
	accountName string,
	accountNo string,
	bankBranch string,
	bankCode string,
	bankName string,
	billPayload string,
	billRes string,
) error {
	q := `
		UPDATE payments
		SET updated_at = now(),
			appotapay_bill_code = $2,
			appotapay_account_no = $3,
			appotapay_account_name = $4,
			appotapay_bank_code = $5,
			appotapay_bank_name = $6,
			appotapay_bank_branch = $7,
			post_data = $8,
			response_data = $9
		WHERE id = $1
	`

	_, err := tx.ExecContext(
		ctx,
		q,
		paymentId,
		billCode,
		accountNo,
		accountName,
		bankCode,
		bankName,
		bankBranch,
		billPayload,
		billRes,
	)

	return err
}
