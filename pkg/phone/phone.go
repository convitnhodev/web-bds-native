package phone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deeincom/deeincom/pkg/models"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// https://esms.vn/eSMS.vn_TailieuAPI.pdf
var ESMS_APIKEY string
var ESMS_SECRET string
var Brandname = "Baotrixemay"

func SendVerifyPhone(user *models.User) error {
	content := fmt.Sprintf("%s la ma xac minh dang ky %s cua ban", user.PhoneToken, Brandname)
	phone := user.Phone

	return SendSMS(phone, content)
}

func SendResetPwdPhone(phone string, token string) error {
	// content sms là cố định, ko thể thay text !
	content := fmt.Sprintf("%s la ma xac minh dang ky %s cua ban", token, Brandname)

	return SendSMS(phone, content)
}

func SendCheckoutCode(phone string, serectCode string) error {
	content := fmt.Sprintf("%s la ma xac minh mua tai san tai %s cua ban", serectCode, Brandname)

	return SendSMS(phone, content)
}

func SendSMS(phone string, content string) error {
	if ESMS_APIKEY == "" {
		return errors.New("ESMS_APIKEY=?")
	}

	if ESMS_SECRET == "" {
		return errors.New("ESMS_SECRET=?")
	}

	fmt.Println("DEBUG", "phone.SendSMS")

	values, err := json.Marshal(map[string]string{
		"ApiKey":    ESMS_APIKEY,
		"SecretKey": ESMS_SECRET,
		"Brandname": Brandname,
		"Content":   content,
		"Phone":     phone,
		"SmsType":   "2",
		// "RequestId": fmt.Sprintf("phone:sendSMS:user:%d", user.ID),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://rest.esms.vn/MainService.svc/json/SendMultipleMessage_V4_post_json/",
		bytes.NewBuffer([]byte(values)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
