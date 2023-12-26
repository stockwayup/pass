package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stockwayup/pass/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command execution failed")
		os.Exit(1)
	}
}
