package service

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stockwayup/pass/conf"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/transport/event"
)

func newPassSvc() *Password {
	cfg := &conf.Config{}
	cfg.Password.Time = 1
	cfg.Password.Memory = 65536
	cfg.Password.Threads = 1
	cfg.Password.KeyLen = 32
	return NewPasswordSvc(cfg)
}

func TestValidator_run_Valid(t *testing.T) {

	svc := newPassSvc()
	v := Validator{password: svc}

	hash, salt, err := svc.HashPassword(context.Background(), []byte("secret"))
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	in := event.Validate{Input: []byte("secret"), Password: hash, Salt: salt}
	data, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var captured *nats.Msg
	oldRespondMsg := respondMsg
	respondMsg = func(_ *nats.Msg, reply *nats.Msg) error { captured = reply; return nil }
	defer func() { respondMsg = oldRespondMsg }()

	msg := &nats.Msg{Header: nats.Header{}}
	msg.Header.Set("id", "v1")
	msg.Data = data

	if err := v.run(context.Background(), msg); err != nil {
		t.Fatalf("run: %v", err)
	}
	if captured == nil {
		t.Fatalf("no reply")
	}
	if got := captured.Header.Get("type"); got != dictionary.TypeValid {
		t.Errorf("type=%q", got)
	}
}

func TestValidator_run_Invalid(t *testing.T) {

	svc := newPassSvc()
	v := Validator{password: svc}

	hash, salt, err := svc.HashPassword(context.Background(), []byte("secret"))
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	in := event.Validate{Input: []byte("wrong"), Password: hash, Salt: salt}
	data, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var captured *nats.Msg
	oldRespondMsg := respondMsg
	respondMsg = func(_ *nats.Msg, reply *nats.Msg) error { captured = reply; return nil }
	defer func() { respondMsg = oldRespondMsg }()

	msg := &nats.Msg{Header: nats.Header{}}
	msg.Header.Set("id", "v2")
	msg.Data = data

	if err := v.run(context.Background(), msg); err != nil {
		t.Fatalf("run: %v", err)
	}
	if captured == nil {
		t.Fatalf("no reply")
	}
	if got := captured.Header.Get("type"); got != dictionary.TypeInvalid {
		t.Errorf("type=%q", got)
	}
}

func TestValidator_run_UnmarshalError(t *testing.T) {

	v := Validator{password: newPassSvc()}

	var dataSent []byte
	oldRespond := respond
	respond = func(_ *nats.Msg, data []byte) error { dataSent = append([]byte(nil), data...); return nil }
	defer func() { respond = oldRespond }()

	msg := &nats.Msg{Header: nats.Header{}}
	msg.Header.Set("id", "v3")
	msg.Data = []byte("bad")

	if err := v.run(context.Background(), msg); err != nil {
		t.Fatalf("run: %v", err)
	}
	if string(dataSent) != dictionary.TypeValidatedError {
		t.Errorf("expected %q, got %q", dictionary.TypeValidatedError, string(dataSent))
	}
}
