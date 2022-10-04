module github.com/stockwayup/pass

go 1.19

replace github.com/isayme/go-amqp-reconnect v0.0.0-20210303120416-fc811b0bcda2 => github.com/soulgarden/go-amqp-reconnect v0.0.0-20221004062723-736d34abd6d3

require (
	github.com/isayme/go-amqp-reconnect v0.0.0-20210303120416-fc811b0bcda2
	github.com/jinzhu/configor v1.2.1
	github.com/rs/zerolog v1.27.0
	github.com/satori/go.uuid v1.2.0
	github.com/soulgarden/rmq-pubsub v0.0.3
	github.com/spf13/cobra v1.5.0
	github.com/streadway/amqp v1.0.0
	github.com/tinylib/msgp v1.1.6
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
