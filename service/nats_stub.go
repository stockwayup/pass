package service

import "github.com/nats-io/nats.go"

// respond and respondMsg provide a test seam for NATS replies.
// They are overridden in unit tests to avoid real NATS calls.
var respond = func(msg *nats.Msg, data []byte) error { return msg.Respond(data) }

var respondMsg = func(msg *nats.Msg, reply *nats.Msg) error { return msg.RespondMsg(reply) }
