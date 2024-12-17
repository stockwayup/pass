package service

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/transport/event"
)

type Generator struct {
	passwordSvc *Password
}

func NewGenerator(passwordSvc *Password) *Generator {
	return &Generator{passwordSvc: passwordSvc}
}

func (s *Generator) Process(
	ctx context.Context,
	delivery <-chan *nats.Msg,
) error {
	for {
		select {
		case msg, ok := <-delivery:
			if !ok {
				return dictionary.ErrDeliveryChannelClosed
			}

			msgID := msg.Header.Get("id")

			ctx := zerolog.Ctx(ctx).With().
				Str("id", msgID).
				Logger().
				WithContext(context.WithValue(ctx, dictionary.ID, msgID))

			in := event.Generate{}

			_, err := in.UnmarshalMsg(msg.Data)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("unmarshal msg")

				if err := msg.Respond([]byte(dictionary.TypeGeneratedError)); err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

					return err
				}

				return err
			}

			hash, salt, err := s.passwordSvc.HashPassword(ctx, in.Password)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("hash password")

				if err := msg.Respond([]byte(dictionary.TypeGeneratedError)); err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

					return err
				}
			}

			out := event.Generated{Hash: hash, Salt: salt}

			body, err := out.MarshalMsg(nil)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("marshal msg")

				if err := msg.Respond([]byte(dictionary.TypeGeneratedError)); err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

					return err
				}
			}

			reply := nats.NewMsg("")

			reply.Header.Set("id", msgID)
			reply.Header.Set("type", dictionary.TypeGenerated)
			reply.Data = body

			if err := msg.RespondMsg(reply); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
