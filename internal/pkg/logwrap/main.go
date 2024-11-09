// Package logwrap is wrapper for zerolog.
package logwrap

import (
	"fmt"
	"os"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/k0st1a/gophkeeper/internal/pkg/traceid"
)

func New(level string) error {
	// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
	//nolint // need here
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			traceidFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
		FieldsExclude: []string{
			traceidFieldName,
		},
	})

	// https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
	//nolint // need here
	log.Logger = log.With().Caller().Logger()

	// https://github.com/rs/zerolog?tab=readme-ov-file#contextcontext-integration
	//nolint // need here
	log.Logger = log.Logger.Hook(tracingHook{})

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

	// for the case when the logger does not exist in the context
	//nolint // need here
	zerolog.DefaultContextLogger = &log.Logger

	return nil
}

func NewFile(level string, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		return fmt.Errorf("error of open file for log:%w", err)
	}

	// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
	//nolint // need here
	log.Logger = log.Output(zerolog.ConsoleWriter{
		NoColor: true,
		Out:     file,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			traceidFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
		FieldsExclude: []string{
			traceidFieldName,
		},
	})

	// https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
	//nolint // need here
	log.Logger = log.With().Caller().Logger()

	// https://github.com/rs/zerolog?tab=readme-ov-file#contextcontext-integration
	//nolint // need here
	log.Logger = log.Logger.Hook(tracingHook{})

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

	// for the case when the logger does not exist in the context
	//nolint // need here
	zerolog.DefaultContextLogger = &log.Logger

	return nil
}

type tracingHook struct{}

const traceidFieldName = "traceid"

func (h tracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	e.Str(traceidFieldName, traceid.Get(ctx))
}
