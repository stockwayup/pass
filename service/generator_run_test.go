package service

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stockwayup/pass/conf"
	"github.com/stockwayup/pass/dictionary"
	"github.com/stockwayup/pass/transport/event"
)

func TestGenerator_run_Success(t *testing.T) {
	// stub responder to capture reply and ensure no error path
	var captured *nats.Msg
	var errRespondCalled bool
	oldRespondMsg := respondMsg
	oldRespond := respond
	respondMsg = func(_ *nats.Msg, reply *nats.Msg) error { captured = reply; return nil }
	respond = func(_ *nats.Msg, _ []byte) error { errRespondCalled = true; return nil }
	defer func() { respondMsg = oldRespondMsg; respond = oldRespond }()

	cfg := &conf.Config{}
	cfg.Password.Time = 1
	cfg.Password.Memory = 65536
	cfg.Password.Threads = 1
	cfg.Password.KeyLen = 32
	g := Generator{passwordSvc: NewPasswordSvc(cfg)}

	in := event.Generate{Password: []byte("secret")}
	data, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	msg := &nats.Msg{Header: nats.Header{}}
	msg.Header.Set("id", "abc")
	msg.Data = data

	if err := g.run(context.Background(), msg); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if captured == nil {
		t.Fatalf("no reply captured")
	}
	if errRespondCalled {
		t.Fatalf("unexpected error respond path")
	}
	if got := captured.Header.Get("id"); got != "abc" {
		t.Errorf("id header = %q", got)
	}
	if got := captured.Header.Get("type"); got != dictionary.TypeGenerated {
		t.Errorf("type header = %q", got)
	}

	var out event.Generated
	if _, err := out.UnmarshalMsg(captured.Data); err != nil {
		t.Fatalf("unmarshal reply: %v", err)
	}
	if len(out.Hash) == 0 || len(out.Salt) == 0 {
		t.Errorf("empty hash/salt in reply")
	}
}

func TestGenerator_run_UnmarshalError(t *testing.T) {
	// stub simple respond to capture error type
	var dataSent []byte
	oldRespond := respond
	respond = func(_ *nats.Msg, data []byte) error { dataSent = append([]byte(nil), data...); return nil }
	defer func() { respond = oldRespond }()

	cfg := &conf.Config{}
	cfg.Password.Time = 1
	cfg.Password.Memory = 65536
	cfg.Password.Threads = 1
	cfg.Password.KeyLen = 32
	g := Generator{passwordSvc: NewPasswordSvc(cfg)}

	msg := &nats.Msg{Header: nats.Header{}}
	msg.Header.Set("id", "abc")
	msg.Data = []byte("not-msgpack")

	if err := g.run(context.Background(), msg); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if string(dataSent) != dictionary.TypeGeneratedError {
		t.Errorf("expected %q, got %q", dictionary.TypeGeneratedError, string(dataSent))
	}
}
