package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure logger level and output format
func Configure(level string, pretty bool, sentryDSN string) {
	zerolog.TimeFieldFormat = ""
	zerolog.DurationFieldInteger = true
	l := zerolog.InfoLevel
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	}
	zerolog.SetGlobalLevel(l)
	var w io.Writer = os.Stdout
	if pretty {
		w = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	if sentryDSN != "" {
		w = zerolog.MultiLevelWriter(w, SentryWriter(sentryDSN))
	}
	log.Logger = zerolog.New(w).With().Timestamp().Logger()
}
