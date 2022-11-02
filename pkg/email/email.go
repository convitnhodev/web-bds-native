package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deeincom/deeincom/pkg/models"
)

var PostmarkApiToken string

var client = &http.Client{
	Timeout: 45 * time.Second,
}

func SendVerifyEmail(user *models.User) error {
	if PostmarkApiToken == "" {
		return errors.New("PostmarkApiToken=?")
	}

	// dùng text cho gọn
	msg := fmt.Sprintf(`Hey %s %s!\n
	A verification code was required to verify it's your email. \n
	Verification code: %s\n
	You can also follow this link to verify your email address https://deein.com/verify-email?code=%s&iat=%d\n
	If you didn't ask to verify this address, you can ignore this email.\n
	Thanks,\n
	The Deein Team\n`, user.FirstName, user.LastName, user.EmailToken, user.EmailToken, time.Now().UTC().Unix())

	body, err := json.Marshal(map[string]string{
		"From":          "",
		"To":            user.Email,
		"Subject":       "[Deein] Please verify your email",
		"TextBody":      msg,
		"MessageStream": "outbound",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-Postmark-Server-Token", PostmarkApiToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
