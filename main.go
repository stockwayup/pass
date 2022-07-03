package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stockwayup/pass/cmd"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command execution failed")
		os.Exit(1)
	}
}
