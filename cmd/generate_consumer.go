package cmd

import (
	"os"

	"github.com/stockwayup/pass/conf"
	"github.com/stockwayup/pass/service"
	"github.com/stockwayup/pass/storage/rmq"

	"github.com/rs/zerolog"
	pubsub "github.com/soulgarden/rmq-pubsub"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewGenerateConsumerCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "generate_consume",
		Short: "Run generate consumer",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			cfg := conf.New()

			logger := zerolog.New(os.Stdout).With().Caller().Logger()

			if cfg.DebugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			cmdManager := service.NewManager(&logger)

			ctx, _ := cmdManager.ListenSignal()

			ctx = logger.WithContext(ctx)

			g, ctx := errgroup.WithContext(ctx)

			rmqDialer := rmq.NewDialer(cfg, &logger)
			rmqConn, err := rmqDialer.Dial()
			if err != nil {
				logger.Err(err).Msg("rabbitmq failed to establish connection")
				os.Exit(1)
			}

			defer rmqConn.Close()

			pub := pubsub.NewPub(
				rmqConn,
				cfg.RMQ.Queues.GenerateOut,
				pubsub.NewRmq(rmqConn, cfg.RMQ.Queues.GenerateOut, &logger),
				&logger,
			)

			sub := pubsub.NewSub(
				rmqConn,
				service.NewGenerator(service.NewPasswordSvc(cfg), pub),
				pubsub.NewRmq(rmqConn, cfg.RMQ.Queues.GenerateIn, &logger),
				cfg.RMQ.Queues.GenerateIn,
				&logger,
			)

			g.Go(func() error {
				err := pub.StartPublisher(ctx)

				logger.Err(err).Msg("start publisher")

				return err
			})

			g.Go(func() error {
				err := sub.StartConsumer(ctx)

				logger.Err(err).Msg("start subscriber")

				return err
			})

			logger.Err(g.Wait()).Msg("wait goroutines")
		},
	}
}
