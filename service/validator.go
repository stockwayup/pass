package service

import (
	"context"

	"github.com/rs/zerolog"
	pubsub "github.com/soulgarden/rmq-pubsub"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/storage/rmq/event"
	"github.com/streadway/amqp"
)

type Validator struct {
	passwordSvc *Password
	pub         *pubsub.Pub
}

func NewValidator(passwordSvc *Password, pub *pubsub.Pub) *Validator {
	return &Validator{passwordSvc: passwordSvc, pub: pub}
}

func (s *Validator) Process(
	ctx context.Context,
	delivery <-chan amqp.Delivery,
) {
	for {
		select {
		case msg := <-delivery:
			ctx := zerolog.Ctx(ctx).With().
				Str("id", msg.MessageId).
				Logger().
				WithContext(context.WithValue(ctx, dictionary.ID, msg.MessageId))

			zerolog.Ctx(ctx).Err(msg.Ack(false)).Str("id", msg.MessageId).Msg("ack")

			in := event.Validate{}

			_, err := in.UnmarshalMsg(msg.Body)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("unmarshal msg")

				s.pub.Publish(event.NewAMQPValidatedMsg(msg.MessageId, dictionary.TypeValidatedError, []byte{}))

				return
			}

			isValid, err := s.passwordSvc.IsValid(in.Input, in.Password, in.Salt)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("is valid")

				s.pub.Publish(event.NewAMQPValidatedMsg(msg.MessageId, dictionary.TypeValidatedError, []byte{}))

				return
			}

			out := event.Validated{IsValid: isValid}

			body, err := out.MarshalMsg(nil)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("marshal msg")

				s.pub.Publish(event.NewAMQPValidatedMsg(msg.MessageId, dictionary.TypeValidatedError, []byte{}))

				return
			}

			s.pub.Publish(event.NewAMQPValidatedMsg(msg.MessageId, dictionary.TypeValidated, body))

		case <-ctx.Done():
			return
		}
	}
}
