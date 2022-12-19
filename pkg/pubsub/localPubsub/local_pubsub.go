package localPubsub

import (
	"context"
	"github.com/deeincom/deeincom/pkg/pubsub"
	"sync"
)

type localPubsub struct {
	messageQueue chan *pubsub.Message
	mapChannel   map[pubsub.Topic][]chan *pubsub.Message
	locker       *sync.RWMutex
}

func NewPubSub() *localPubsub {
	pb := &localPubsub{
		make(chan *pubsub.Message, 10000),
		make(map[pubsub.Topic][]chan *pubsub.Message),
		new(sync.RWMutex),
	}
	pb.run()
	return pb

}

func (ps *localPubsub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	data.SetChannel(topic)

	go func() {
		ps.messageQueue <- data

	}()
	return nil
}

func (ps *localPubsub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)
	ps.locker.Lock()

	if val, ok := ps.mapChannel[topic]; ok {
		val = append(ps.mapChannel[topic], c)
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *pubsub.Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					chans = append(chans[:i], chans[i+1:]...)
					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()
					break
				}
			}
		}
	}
}

func (ps *localPubsub) run() error {
	go func() {
		for {
			mess := <-ps.messageQueue

			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				for i := range subs {
					go func(c chan *pubsub.Message) {
						c <- mess
					}(subs[i])
				}
			}
		}
	}()
	return nil
}