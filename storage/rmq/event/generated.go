package event

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stockwayup/pass/dictionary"
)

//go:generate msgp

type Generated struct {
	Hash []byte `msg:"hash"`
	Salt []byte `msg:"salt"`
}

func NewAMQPGeneratedMsg(id string, t string, body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType:  dictionary.QueueMessageContentType,
		Type:         t,
		MessageId:    id,
		DeliveryMode: amqp.Transient,
		Body:         body,
		Timestamp:    time.Now(),
		Expiration:   dictionary.TTL,
	}
}
