package broker

import (
	"github.com/pu4mane/goChat/internal/app/model"
)

type MessageBroker interface {
	Subscribe(subject string, callback func(msg *model.Message)) (interface{}, error)
	Publish(subject string, message *model.Message) error
	Unsubscribe(subscription interface{})
}
