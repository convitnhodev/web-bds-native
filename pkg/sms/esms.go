package sms

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type Payload struct {
	APIKey      string `json:"ApiKey"`
	Content     string `json:"Content"`
	Phone       string `json:"Phone"`
	SecretKey   string `json:"SecretKey"`
	IsUnicode   string `json:"IsUnicode"`
	BrandName   string `json:"Brandname"`
	SmsType     string `json:"SmsType"`
	RequestID   string `json:"RequestId"`
	CallbackURL string `json:"CallbackUrl"`
	CampaignID  string `json:"campaignid"`
}

type SMSSender struct {
	apiKey    string
	secretKey string
	brand     string
	c         *http.Client
}

type SentResp struct {
	ID              string `json:"-"`
	CodeResult      string `json:"CodeResult"`
	CountRegenerate int    `json:"CountRegenerate"`
	SmsID           string `json:"SMSID"`
}

func NewSMSClient(apiKey string, secretKey string, brand string) *SMSSender {
	return &SMSSender{c: &http.Client{}, apiKey: apiKey, secretKey: secretKey, brand: brand}
}

func (s *SMSSender) Sent(pn string, content string) (*SentResp, error) {
	url := "http://rest.esms.vn/MainService.svc/json/SendMultipleMessage_V4_post_json"
	id := uuid.New().String()
	payload := Payload{
		APIKey:      s.apiKey,
		Content:     content,
		Phone:       pn,
		SecretKey:   s.secretKey,
		IsUnicode:   "0",
		BrandName:   s.brand,
		SmsType:     "2",
		RequestID:   id,
		CallbackURL: "",
		CampaignID:  "",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, &buf)
	req.Header.Add("Content-Type", "application/json")
	res, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var rd SentResp
	if err := json.Unmarshal(body, &rd); err != nil {
		return nil, err
	}
	rd.ID = id
	return &rd, nil
}
