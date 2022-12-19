package pubsub

import (
	"fmt"
	"time"
)

type Message struct {
	id        string
	chanel    Topic
	data      interface{}
	createdAt time.Time
}

func NewMessage(data interface{}) *Message {
	now := time.Now().UTC()
	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %s", evt.chanel)
}

func (evt *Message) Channel() Topic {
	return evt.chanel
}

func (evt *Message) SetChannel(channel Topic) {
	evt.chanel = channel
}

func (evt *Message) Data() interface{} {
	return evt.data
}