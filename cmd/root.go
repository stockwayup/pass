package cmd

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

func NewRootCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "stockwayup",
		Short: "Manage passwords",
		Args:  cobra.ExactArgs(1),
	}
}

func Execute() error {
	root := NewRootCMD()
	root.AddCommand(NewGenerateConsumerCMD())
	root.AddCommand(NewValidateConsumerCMD())

	root.InitDefaultHelpFlag()

	if err := root.Execute(); err != nil {
		log.Err(err).Msg("command execution failed")

		return err
	}

	return nil
}
