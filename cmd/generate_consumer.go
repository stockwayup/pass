package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/stockwayup/pass/transport"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stockwayup/pass/conf"
	"github.com/stockwayup/pass/service"
	"golang.org/x/sync/errgroup"
)

const generatorWorkerName = "swup.pass.generation"

func NewGenerateConsumerCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "generate_consume",
		Short: "Run generate consumer",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := conf.New()

			logger := zerolog.New(os.Stdout).With().Caller().Logger()

			if cfg.DebugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			cmdManager := service.NewManager(&logger)

			ctx, _ := cmdManager.ListenSignal()

			ctx = logger.WithContext(ctx)

			g, ctx := errgroup.WithContext(ctx)

			const serviceName = "pass.generate"

			natsConn, err := transport.NewConnection(cfg, serviceName)
			if err != nil {
				logger.Err(err).Msg("nats failed to establish connection")
				os.Exit(1)
			}

			defer natsConn.Close()

			mch := make(chan *nats.Msg, natsConn.Opts.SubChanLen)

			sub, err := natsConn.ChanQueueSubscribe(cfg.Nats.Queues.Generation, generatorWorkerName, mch)
			if err != nil {
				logger.Err(err).Msg("nats failed to subscribe")
				os.Exit(1)
			}

			defer sub.Unsubscribe()

			g.Go(func() error {
				err := service.NewGenerator(service.NewPasswordSvc(cfg)).Process(ctx, mch)

				logger.Err(err).Msg("generator process")

				return err
			})

			err = g.Wait()
			if err == nil || errors.Is(context.Canceled, err) {
				return nil
			}

			logger.Err(err).Msg("wait goroutines")

			return err
		},
	}
}
