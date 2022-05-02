package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func initLogging() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
	log = zerolog.New(output).With().Timestamp().Logger()
	if verboseFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
