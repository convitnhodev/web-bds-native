package pubsub

import "context"

type Topic string

type Pubsub interface {
	Publish(ctx context.Context, topic Topic, data *Message) error
	Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, close func())
}