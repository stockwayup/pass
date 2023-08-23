package event

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	uuid "github.com/satori/go.uuid"
	"github.com/stockwayup/pass/dictionary"
)

//go:generate msgp

type Generate struct {
	Password []byte `msg:"password"`
}

func NewGenerateMsg(e Generate) (p amqp.Publishing, _ error) {
	body, err := e.MarshalMsg(nil)
	if err != nil {
		return p, err
	}

	return amqp.Publishing{
		ContentType:  dictionary.QueueMessageContentType,
		MessageId:    uuid.NewV4().String(),
		DeliveryMode: amqp.Transient,
		Body:         body,
		Timestamp:    time.Now(),
	}, nil
}
