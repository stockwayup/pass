package event

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stockwayup/pass/dictionary"
)

//go:generate msgp

type Validated struct {
	IsValid bool `msg:"is_valid"`
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
