package appotapay

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Payment
type APTPaymentPayload struct {
	Amount        int    `json:"amount"`
	OrderId       string `json:"orderId"`
	OrderInfo     string `json:"orderInfo"`
	BankCode      string `json:"bankCode"`
	PaymentMethod string `json:"paymentMethod"`
	ClientIP      string `json:"clientIp"`
	ExtraData     string `json:"extraData"`
	NotifyUrl     string `json:"notifyUrl"`
	RedirectUrl   string `json:"redirectUrl"`
	Signature     string `json:"signature"`
}

func (m *APTPaymentPayload) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&bankCode=%s&clientIp=%s&extraData=%s&notifyUrl=%s&orderId=%s&orderInfo=%s&paymentMethod=%s&redirectUrl=%s",
		m.Amount,
		m.BankCode,
		m.ClientIP,
		m.ExtraData,
		m.NotifyUrl,
		m.OrderId,
		m.OrderInfo,
		m.PaymentMethod,
		m.RedirectUrl,
	)
}

type APTPaymentResponse struct {
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	OrderId    string `json:"orderId"`
	Ammount    int    `json:"amount"`
	PaymentUrl string `json:"paymentUrl"`
	Signature  string `json:"signature"`
}

func (m *APTPaymentResponse) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&errorCode=%d&message=%s&orderId=%s&paymentUrl=%s",
		m.Ammount,
		m.ErrorCode,
		m.Message,
		m.OrderId,
		m.PaymentUrl,
	)
}

// Refund Payment
type APTRefundData struct {
	AppotapayTransId string `json:"appotapayTransId"`
	RefundId         string `json:"refundId"`
	RefundOriginalId string `json:"refundOriginalId"`
	Ammount          int    `json:"amount"`
	Reason           string `json:"reason"`
	Status           string `json:"status"`
	TransactionTs    int    `json:"transactionTs"`
}

type APTRefundResponse struct {
	ErrorCode int           `json:"errorCode"`
	Message   string        `json:"message"`
	Data      APTRefundData `json:"data"`
	Signature string        `json:"signature"`
}

func (m *APTRefundResponse) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&appotapayTransId=%s&errorCode=%d&reason=%s&refundId=%s&refundOriginalId=%s&status=%s&transactionTs=%d",
		m.Data.Ammount,
		m.Data.AppotapayTransId,
		m.ErrorCode,
		m.Data.Reason,
		m.Data.RefundId,
		m.Data.RefundOriginalId,
		m.Data.Status,
		m.Data.TransactionTs,
	)
}

type APTRefundPayload struct {
	RefundId         string `json:"refundId"`
	AppotapayTransId string `json:"appotapayTransId"`
	Amount           int    `json:"amount"`
	Reason           string `json:"reason"`
	Signature        string `json:"signature"`
}

func (m *APTRefundPayload) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&appotapayTransId=%s&reason=%s&refundId=%s",
		m.Amount,
		m.AppotapayTransId,
		m.Reason,
		m.RefundId,
	)
}

type APTPaymentRecipition struct {
	ErrorCode        int    `json:"errorCode"`
	Message          string `json:"message"`
	PartnerCode      string `json:"partnerCode"`
	ApiKey           string `json:"apiKey"`
	Amount           int    `json:"amount"`
	Currency         string `json:"currency"`
	OrderId          string `json:"orderId"`
	BankCode         string `json:"bankCode"`
	PaymentMethod    string `json:"paymentMethod"`
	PaymentType      string `json:"paymentType"`
	AppotapayTransId string `json:"appotapayTransId"`
	TransactionTs    int    `json:"transactionTs"`
	ExtraData        string `json:"extraData"`
	Signature        string `json:"signature"`
}

func (o *APTPaymentRecipition) ParseOrderId(aptPaymentHost string) string {
	strSlice := strings.Split(o.OrderId, "-")

	if strings.Contains(aptPaymentHost, "payment.dev.appotapay.com") && len(strSlice) > 1 {
		return strSlice[1]
	}
	return o.OrderId
}

func (m *APTPaymentRecipition) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&apiKey=%s&appotapayTransId=%s&bankCode=%s&currency=%s&errorCode=%d&extraData=%s&message=%s"+
			"&orderId=%s&partnerCode=%s&paymentMethod=%s&paymentType=%s&transactionTs=%d",
		m.Amount,
		m.ApiKey,
		m.AppotapayTransId,
		m.BankCode,
		m.Currency,
		m.ErrorCode,
		m.ExtraData,
		m.Message,
		m.OrderId,
		m.PartnerCode,
		m.PaymentMethod,
		m.PaymentType,
		m.TransactionTs,
	)
}

type APTBillPayload struct {
	BillCode         string `json:"billCode"`
	BillInfo         string `json:"billInfo"`
	BillExpiryTime   int64  `json:"billExpiryTime"`
	CustomerName     string `json:"customerName"`
	ServiceCode      string `json:"serviceCode"`
	Amount           int    `json:"amount"`
	PaymentCondition string `json:"paymentCondition"`
	BankCode         string `json:"bankCode"`
	NotifyUrl        string `json:"notifyUrl"`
	ExtraData        string `json:"extraData"`
	Signature        string `json:"signature"`
}

func (m *APTBillPayload) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&bankCode=%s&billCode=%s&billExpiryTime=%d&billInfo=%s&customerName=%s&extraData=%s&notifyUrl=%s&paymentCondition=%s&serviceCode=%s",
		m.Amount,
		m.BankCode,
		m.BillCode,
		m.BillExpiryTime,
		m.BillInfo,
		m.CustomerName,
		m.ExtraData,
		m.NotifyUrl,
		m.PaymentCondition,
		m.ServiceCode,
	)
}

type APTBillAccountBank struct {
	BankCode    string `json:"bankCode"`
	BankName    string `json:"bankName"`
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	BankBranch  string `json:"bankBranch"`
}

type APTBillPayment struct {
	BankAccounts []APTBillAccountBank `json:"bankAccounts"`
}

type APTBillResponse struct {
	ErrorCode int            `json:"errorCode"`
	Message   string         `json:"message"`
	BillCode  string         `json:"billCode"`
	Payment   APTBillPayment `json:"payment"`
	Signature string         `json:"signature"`
}

func (m *APTBillResponse) GetPayloadUrl() string {
	paymentJson, _ := json.Marshal(m.Payment)
	return fmt.Sprintf(
		"billCode=%s&errorCode=%d&payment=%s",
		m.BillCode,
		m.ErrorCode,
		string(paymentJson),
	)
}

type APTBillRecipition struct {
	BillCode          string `json:"message"`
	PartnerCode       string `json:"partnerCode"`
	ApiKey            string `json:"apiKey"`
	Amount            int    `json:"amount"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankAccountName   string `json:"bankAccountName"`
	BankCode          string `json:"bankCode"`
	RequestTime       int    `json:"requestTime"`
	TransactionTime   int    `json:"transactionTime"`
	TransactionId     string `json:"appotapayTransId"`
	Version           string `json:"version"`
	Memo              string `json:"memo"`
	ExtraData         string `json:"extraData"`
	Signature         string `json:"signature"`
}

func (m *APTBillRecipition) GetPayloadUrl() string {
	return fmt.Sprintf(
		"amount=%d&apiKey=%s&bankAccountName=%s&bankAccountNumber=%s&bankCode=%s&billCode=%s&extraData=%s&memo=%s&partnerCode=%s&requestTime=%d&transactionId=%s&transactionTime=%d&version=%s",
		m.Amount,
		m.ApiKey,
		m.BankAccountName,
		m.BankAccountNumber,
		m.BankCode,
		m.BillCode,
		m.ExtraData,
		m.Memo,
		m.PartnerCode,
		m.RequestTime,
		m.TransactionId,
		m.TransactionTime,
		m.Version,
	)
}

func (m *APTBillRecipition) ParseOrderId(aptEbillHost string) string {
	strSlice := strings.Split(m.BillCode, "-")

	if strings.Contains(aptEbillHost, "payment.dev.appotapay.com") && len(strSlice) > 1 {
		return strSlice[1]
	}
	return m.BillCode
}
