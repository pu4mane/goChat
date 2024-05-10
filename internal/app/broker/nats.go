package broker

import (
	"github.com/nats-io/nats.go"
	"github.com/pu4mane/goChat/internal/app/model"
)

type NATS struct {
	conn *nats.Conn
}

func NewNATS(URL string) (*NATS, error) {
	conn, err := nats.Connect(URL)
	if err != nil {
		return nil, err
	}

	return &NATS{
		conn: conn,
	}, nil
}

func (ns *NATS) Subscribe(subject string, callback func(msg *model.Message)) (interface{}, error) {
	subscription, err := ns.conn.Subscribe(subject, func(msg *nats.Msg) {
		callback(&model.Message{Text: string(msg.Data)})
	})
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (ns *NATS) Unsubscribe(subscription interface{}) {
	if sub, ok := subscription.(*nats.Subscription); ok {
		sub.Unsubscribe()
	}
}

func (ns *NATS) Publish(subject string, message *model.Message) error {
	return ns.conn.Publish(subject, []byte(message.Text))
}
