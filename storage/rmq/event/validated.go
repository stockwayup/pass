package event

import (
	"time"

	"github.com/stockwayup/pass/dictionary"
	"github.com/streadway/amqp"
)

//go:generate msgp

type Validated struct {
	IsValid bool `msgp:"is_valid"`
}

func NewAMQPValidatedMsg(id string, t string, body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType:  dictionary.QueueMessageContentType,
		MessageId:    id,
		Type:         t,
		DeliveryMode: amqp.Transient,
		Body:         body,
		Timestamp:    time.Now(),
		Expiration:   dictionary.TTL,
	}
}
