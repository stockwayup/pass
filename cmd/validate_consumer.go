package cmd

import (
	"github.com/nats-io/nats.go"
	"github.com/stockwayup/pass/transport"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stockwayup/pass/conf"
	"github.com/stockwayup/pass/service"
	"golang.org/x/sync/errgroup"
)

const validatorWorkerName = "swup.pass.validation"

func NewValidateConsumerCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "validate_consume",
		Short: "Run validate consumer",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			var (
				cfg    = conf.New()
				logger = zerolog.New(os.Stdout).With().Caller().Logger()
			)

			if cfg.DebugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			var (
				cmdManager = service.NewManager(&logger)
				ctx, _     = cmdManager.ListenSignal()
			)

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

			sub, err := natsConn.ChanQueueSubscribe(cfg.Nats.Queues.Validation, validatorWorkerName, mch)
			if err != nil {
				logger.Err(err).Msg("nats failed to subscribe")
				os.Exit(1)
			}

			defer sub.Unsubscribe()

			g.Go(func() error {
				err := service.NewValidator(service.NewPasswordSvc(cfg)).Process(ctx, mch)

				logger.Err(err).Msg("validator process")

				return err
			})

			logger.Err(g.Wait()).Msg("wait goroutines")
		},
	}
}
