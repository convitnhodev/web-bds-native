package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/deeincom/deeincom/pkg/helper"
)

type Telegram struct {
	Token  string
	ChatID string
}

type payload struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (o *Telegram) Msg(ctx context.Context, s string) error {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", o.Token)

	data := url.Values{}
	data.Set("chat_id", o.ChatID)
	data.Set("text", s)

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := helper.Fetch(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	p := new(payload)
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return err
	}

	if !p.OK {
		return errors.New(p.Description)
	}

	return nil
}
