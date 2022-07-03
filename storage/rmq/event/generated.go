package event

import (
	"time"

	"github.com/stockwayup/pass/dictionary"
	"github.com/streadway/amqp"
)

//go:generate msgp

type Generated struct {
	Hash []byte `msgp:"hash"`
	Salt []byte `msgp:"salt"`
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
