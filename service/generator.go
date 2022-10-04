package service

import (
	"context"

	"github.com/rs/zerolog"
	pubsub "github.com/soulgarden/rmq-pubsub"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/storage/rmq/event"
	"github.com/streadway/amqp"
)

type Generator struct {
	passwordSvc *Password
	pub         pubsub.Pub
}

func NewGenerator(passwordSvc *Password, pub pubsub.Pub) *Generator {
	return &Generator{passwordSvc: passwordSvc, pub: pub}
}

func (s *Generator) Process(
	ctx context.Context,
	delivery <-chan amqp.Delivery,
) error {
	for {
		select {
		case msg, ok := <-delivery:
			if !ok {
				return dictionary.ErrDeliveryChannelClosed
			}

			ctx := zerolog.Ctx(ctx).With().
				Str("id", msg.MessageId).
				Logger().
				WithContext(context.WithValue(ctx, dictionary.ID, msg.MessageId))

			zerolog.Ctx(ctx).Err(msg.Ack(false)).Str("id", msg.MessageId).Msg("ack")

			in := event.Generate{}

			_, err := in.UnmarshalMsg(msg.Body)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("unmarshal msg")

				s.pub.Publish(event.NewAMQPGeneratedMsg(msg.MessageId, dictionary.TypeGeneratedError, []byte{}))

				return err
			}

			hash, salt, err := s.passwordSvc.HashPassword(ctx, in.Password)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("hash password")

				s.pub.Publish(event.NewAMQPGeneratedMsg(msg.MessageId, dictionary.TypeGeneratedError, []byte{}))

				return err
			}

			out := event.Generated{Hash: hash, Salt: salt}

			body, err := out.MarshalMsg(nil)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("marshal msg")

				s.pub.Publish(event.NewAMQPGeneratedMsg(msg.MessageId, dictionary.TypeGeneratedError, []byte{}))

				return err
			}

			s.pub.Publish(event.NewAMQPGeneratedMsg(msg.MessageId, dictionary.TypeGenerated, body))

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
