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

func SendSMS(user *models.User) error {
	if ESMS_APIKEY == "" {
		return errors.New("ESMS_APIKEY=?")
	}

	if ESMS_SECRET == "" {
		return errors.New("ESMS_SECRET=?")
	}

	values, err := json.Marshal(map[string]string{
		"ApiKey":    ESMS_APIKEY,
		"SecretKey": ESMS_SECRET,
		"Brandname": Brandname,
		"Content":   fmt.Sprintf("%s la ma xac minh dang ky %s cua ban", user.PhoneToken, Brandname),
		"Phone":     "",
		"IsUnicode": "0",
		"SmsType":   "2",
		"RequestId": fmt.Sprintf("user:verify:%d", user.ID),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://rest.esms.vn/MainService.svc/json/SendMultipleMessage_V4_post_json/", bytes.NewBuffer([]byte(values)))
	if err != nil {
		return err
	}

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
