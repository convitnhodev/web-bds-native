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

func signSingature(payloadUrl string) string {
	secretKeyBytes := []byte(SecretKey)
	mac := hmac.New(sha256.New, secretKeyBytes)

	mac.Write([]byte(payloadUrl))

	return hex.EncodeToString(mac.Sum(nil))
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

func getAuth() (string, error) {
	t := time.Now().Add(15 * time.Minute)
	jwtPayload := map[string]interface{}{
		"iss":     PartnerCode,
		"jti":     fmt.Sprintf("%s-%d", ApiKey, t.Unix()),
		"api_key": ApiKey,
		"exp":     t.Unix(),
	}

	return signJWT(jwtPayload)
}

// Verfiy IPN Payment
func VerifyIPNPaymentCallback(paymentData APTPaymentRecipition) (string, error) {
	jsonStr := ""
	jsonByte, err := json.Marshal(paymentData)
	if err != nil {
		jsonStr = "{}"
	}
	jsonStr = string(jsonByte)

	signature := signSingature(paymentData.GetPayloadUrl())
	if signature == paymentData.Signature {
		return jsonStr, nil
	}

	return jsonStr, errors.New("Signature không hợp lệ")
}

// Checkout menthod
func Checkout(payload *APTPaymentPayload) (*APTPaymentResponse, error) {
	if PartnerCode == "" {
		return nil, errors.New("PartnerCode=?")
	}

	if ApiKey == "" {
		return nil, errors.New("ApiKey=?")
	}

	if SecretKey == "" {
		return nil, errors.New("SecretKey=?")
	}

	payload.Signature = signSingature(payload.GetPayloadUrl())
	values, err := json.Marshal(payload)

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

	jwtToken, err := getAuth()
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

	var result APTPaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	signature := signSingature(result.GetPayloadUrl())
	if signature != result.Signature {
		return nil, errors.New("Response Appotapay Error")
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Refund menthod
func Refund(payload APTRefundPayload) (*APTRefundResponse, error) {
	if ApiKey == "" {
		return nil, errors.New("ApiKey=?")
	}

	if SecretKey == "" {
		return nil, errors.New("SecretKey=?")
	}

	payload.Signature = signSingature(payload.GetPayloadUrl())
	values, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	httpRequestURL := files.JoinURL(APTPaymentHost, "/api/v1/transaction/refund")
	req, err := http.NewRequest(
		"POST",
		httpRequestURL,
		bytes.NewBuffer([]byte(values)),
	)

	if err != nil {
		return nil, err
	}

	jwtToken, err := getAuth()
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

	var result APTRefundResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	signature := signSingature(result.GetPayloadUrl())
	if signature != result.Signature {
		return nil, errors.New("Response Appotapay Refund Error")
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateBill() {

}
