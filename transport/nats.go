package transport

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stockwayup/pass/conf"
)

func NewConnection(cfg *conf.Config, name string) (*nats.Conn, error) {
	return nats.Connect(cfg.Nats.Host, nats.Name(name), nats.PingInterval(time.Second*30))
}
