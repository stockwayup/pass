package service

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/transport/event"
)

type Validator struct {
	password *Password
}

func NewValidator(passwordSvc *Password) *Validator {
	return &Validator{password: passwordSvc}
}

func (s Validator) Process(
	ctx context.Context,
	delivery <-chan *nats.Msg,
) error {
	for {
		select {
		case msg, ok := <-delivery:
			if !ok {
				return dictionary.ErrDeliveryChannelClosed
			}

			if err := s.run(ctx, msg); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s Validator) run(ctx context.Context, msg *nats.Msg) error {
	msgID := msg.Header.Get("id")

	ctx = zerolog.Ctx(ctx).With().
		Str("id", msgID).
		Logger().
		WithContext(context.WithValue(ctx, dictionary.ID, msgID))

	in := event.Validate{}

	zerolog.Ctx(ctx).Info().Msg("validation request received")
	defer zerolog.Ctx(ctx).Info().Msg("validation request processed")

	if _, err := in.UnmarshalMsg(msg.Data); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("unmarshal msg")

		if err := msg.Respond([]byte(dictionary.TypeValidatedError)); err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

			return err
		}

		return nil
	}

	isValid, err := s.password.IsValid(in.Input, in.Password, in.Salt)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("is valid")

		if err := msg.Respond([]byte(dictionary.TypeValidatedError)); err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

			return err
		}
	}

	resp := dictionary.TypeValid
	if !isValid {
		resp = dictionary.TypeInvalid
	}

	reply := nats.NewMsg("")

	reply.Header.Set("id", msgID)
	reply.Data = []byte(resp)

	if err := msg.RespondMsg(reply); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("nats queue respond")

		return err
	}

	return nil
}
