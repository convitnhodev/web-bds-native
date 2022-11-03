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
var From string = "no-reply@deein.com"

var client = &http.Client{
	Timeout: 45 * time.Second,
}

func SendVerifyEmail(user *models.User) error {
	if PostmarkApiToken == "" {
		return errors.New("PostmarkApiToken=?")
	}

	// dùng text cho gọn
	msg := fmt.Sprintf(`
Chào %s %s!
Để Deein ghi nhận email này là của bạn, chúng tôi cần bạn nhập đúng mã xác thực bên dưới.
Mã xác thực của bạn là: %s
Bạn cũng có thể click vào link này để xác thực email https://deein.com/verify-email?code=%s&iat=%d
Nếu bạn không phải là người yêu cầu xác thực email, bạn có thể bỏ qua bức thư này.
Trân trọng,
Deein`,
		user.FirstName, user.LastName, user.EmailToken, user.EmailToken, time.Now().UTC().Unix())

	body, err := json.Marshal(map[string]string{
		"From":          From,
		"To":            user.Email,
		"Subject":       "[Deein] Vui lòng xác thực email của bạn",
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Postmark-Server-Token", PostmarkApiToken)

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
