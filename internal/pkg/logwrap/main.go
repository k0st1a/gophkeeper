// Package logwrap is wrapper for zerolog.
package logwrap

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(level string) error {
	// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
	//nolint // need here
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
	//nolint // need here
	log.Logger = log.With().Caller().Logger()

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		return fmt.Errorf("unknown log level %v", level)
	}

	return nil
}
