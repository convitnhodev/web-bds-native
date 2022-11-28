package appotapay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// https://esms.vn/eSMS.vn_TailieuAPI.pdf
// PartnerCode: APPOTAPAY
// ApiKey: FJcmF8uj2ISveL5FvvNk4pnp8xrhINz8
// SecretKey: XAonJgy14YhtePEITXhyBS2unjfJLAV3

var ATPHost string
var PartnerCode string
var ApiKey string
var SecretKey string

type ATPPayload struct {
	Amount        int
	OrderId       string
	OrderInfo     string
	BankCode      string
	PaymentMethod string
	ClientIP      string
	ExtraData     string
	NotifyUrl     string
	RedirectUrl   string
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
	mySigningKey := []byte(SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range payload {
		claims[k] = v
	}

	signature, err := token.SignedString(mySigningKey)

	return signature, err
}

func verifyResponse(res ATPResponse) error {
	dataPayload := map[string]interface{}{
		"errorCode":  res.ErrorCode,
		"message":    res.Message,
		"orderId":    res.OrderId,
		"amount":     res.Ammount,
		"paymentUrl": res.PaymentUrl,
	}
	selfSignature, err := signJWT(dataPayload)

	if err != nil {
		return err
	}

	if res.Signature == selfSignature {
		return nil
	}
	return errors.New("Response Appotapay Error")
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
		"orderId":       payload.OrderId,
		"orderInfo":     payload.OrderInfo,
		"bankCode":      payload.BankCode,
		"paymentMethod": payload.PaymentMethod,
		"clientIp":      payload.ClientIP,
		"extraData":     payload.ExtraData,
		"notifyUrl":     payload.NotifyUrl,
		"redirectUrl":   payload.RedirectUrl,
	}

	singSignature, err := signJWT(dataPayload)
	if err != nil {
		return nil, err
	}

	dataPayload["signature"] = singSignature
	values, err := json.Marshal(dataPayload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/orders/payment/bank", ATPHost),
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("typ", "JWT")
	req.Header.Set("alg", "HS256")
	req.Header.Set("cty", "appotapay-api;v=1")
	req.Header.Set("X-APPOTAPAY-AUTH", fmt.Sprintf("Bearer %s", jwtToken))

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
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
