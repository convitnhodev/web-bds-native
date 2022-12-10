package appotapay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/deeincom/deeincom/pkg/files"
	"github.com/golang-jwt/jwt/v4"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

var APTPaymentHost string
var PartnerCode string
var ApiKey string
var SecretKey string

type ATPPayload struct {
	Amount        int64
	OrderId       string
	OrderInfo     string
	BankCode      string
	PaymentMethod string
	ClientIP      string
	ExtraData     string
	NotifyUrl     string
	RedirectUrl   string
}

type PaymentRecipition struct {
	ErrorCode        int    `json:"errorCode"`
	Message          string `json:"message"`
	PartnerCode      string `json:"partnerCode"`
	ApiKey           string `json:"apiKey"`
	Amount           int64  `json:"amount"`
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

func (o *PaymentRecipition) ParseOrderId(APTPaymentHost string) string {
	strSlice := strings.Split(o.OrderId, "-")

	if strings.Contains(APTPaymentHost, "payment.dev.appotapay.com") && len(strSlice) > 1 {
		return strSlice[1]
	}
	return o.OrderId
}

type ATPResponse struct {
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	OrderId    string `json:"orderId"`
	Ammount    int    `json:"amount"`
	PaymentUrl string `json:"paymentUrl"`
	Signature  string `json:"signature"`
}

func signJWT(payload map[string]interface{}) (string, error) {
	secretKeyBytes := []byte(SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range payload {
		claims[k] = v
	}

	signature, err := token.SignedString(secretKeyBytes)

	return signature, err
}

func signSingature(payload map[string]interface{}) string {
	secretKeyBytes := []byte(SecretKey)
	mac := hmac.New(sha256.New, secretKeyBytes)
	// Thứ tự của key được short đừng thay đổi (https://docs.appotapay.com/payment/signature)
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	messages := make([]string, 0, len(payload))
	for _, k := range keys {
		messages = append(messages, fmt.Sprintf("%s=%s", k, fmt.Sprint(payload[k])))
	}

	mac.Write([]byte(strings.Join(messages, "&")))

	return hex.EncodeToString(mac.Sum(nil))
}

func verifyResponse(res ATPResponse) error {
	dataPayload := map[string]interface{}{
		"errorCode":  res.ErrorCode,
		"message":    res.Message,
		"orderId":    res.OrderId,
		"amount":     res.Ammount,
		"paymentUrl": res.PaymentUrl,
	}

	selfSignature := signSingature(dataPayload)

	if res.Signature == selfSignature {
		return nil
	}

	return errors.New("Response Appotapay Error")
}

func VerifyIPNCallback(paymentData PaymentRecipition) (string, error) {
	dataPayload := map[string]interface{}{
		"errorCode":        paymentData.ErrorCode,
		"message":          paymentData.Message,
		"partnerCode":      paymentData.PartnerCode,
		"apiKey":           paymentData.ApiKey,
		"amount":           paymentData.Amount,
		"currency":         paymentData.Currency,
		"orderId":          paymentData.OrderId,
		"bankCode":         paymentData.BankCode,
		"paymentMethod":    paymentData.PaymentMethod,
		"paymentType":      paymentData.PaymentType,
		"appotapayTransId": paymentData.AppotapayTransId,
		"transactionTs":    paymentData.TransactionTs,
		"extraData":        paymentData.ExtraData,
	}
	selfSignature := signSingature(dataPayload)

	jsonStr := ""
	jsonByte, err := json.Marshal(dataPayload)
	if err != nil {
		jsonStr = "{}"
	}
	jsonStr = string(jsonByte)

	if paymentData.Signature == selfSignature {
		return jsonStr, nil
	}

	return jsonStr, errors.New("Signature không hợp lệ")
}

func Checkout(payload *ATPPayload) (*ATPResponse, error) {
	if PartnerCode == "" {
		return nil, errors.New("PartnerCode=?")
	}

	if ApiKey == "" {
		return nil, errors.New("ApiKey=?")
	}

	if SecretKey == "" {
		return nil, errors.New("SecretKey=?")
	}

	dataPayload := map[string]interface{}{
		"amount":        payload.Amount,
		"bankCode":      payload.BankCode,
		"clientIp":      payload.ClientIP,
		"extraData":     payload.ExtraData,
		"notifyUrl":     payload.NotifyUrl,
		"orderId":       payload.OrderId,
		"orderInfo":     payload.OrderInfo,
		"paymentMethod": payload.PaymentMethod,
		"redirectUrl":   payload.RedirectUrl,
	}

	signatureStr := signSingature(dataPayload)
	dataPayload["signature"] = signatureStr
	values, err := json.Marshal(dataPayload)

	if err != nil {
		return nil, err
	}

	httpRequestURL := files.JoinURL(APTPaymentHost, "/api/v1/orders/payment/bank")
	req, err := http.NewRequest(
		"POST",
		httpRequestURL,
		bytes.NewBuffer([]byte(values)),
	)

	if err != nil {
		return nil, err
	}

	t := time.Now().Add(15 * time.Minute)
	jwtPayload := map[string]interface{}{
		"iss":     PartnerCode,
		"jti":     fmt.Sprintf("%s-%d", ApiKey, t.Unix()),
		"api_key": ApiKey,
		"exp":     t.Unix(),
	}
	jwtToken, err := signJWT(jwtPayload)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-APPOTAPAY-AUTH", fmt.Sprintf("Bearer %s", jwtToken))

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	if err != nil {
		return nil, err
	}

	var result ATPResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	err = verifyResponse(result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
