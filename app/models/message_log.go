package models

import (
	"github.com/upper/db/v4"
	"time"
)

var MessageLogTable = "message_logs"

type MessageLog struct {
	ID          int64     `json:"id" db:"id,omitempty"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Content     string    `json:"content" db:"content"`
	RequestID   string    `json:"request_id" db:"request_id"`
	SmsID       string    `json:"sms_id" db:"sms_id"`
	CodeResult  string    `json:"code_result" db:"code_result"`
	SentStatus  string    `json:"sent_status" db:"sent_status"`
	TelcoID     string    `json:"telco_id" db:"telco_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (u *MessageLog) Store(sess db.Session) db.Store {
	return sess.Collection(MessageLogTable)
}

var _ = db.Record(&MessageLog{})
