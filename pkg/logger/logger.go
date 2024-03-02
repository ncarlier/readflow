package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure logger level and output format
func Configure(level, format, sentryDSN string) {
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
	if format == "text" {
		w = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	if sentryDSN != "" {
		w = zerolog.MultiLevelWriter(w, SentryWriter(sentryDSN))
	}
	log.Logger = zerolog.New(w).With().Timestamp().Logger()
}

// Functions redefinition
var (
	// With creates a child logger with the field added to its context.
	With = log.With
	// Debug starts a new message with debug level.
	Debug = log.Debug
	// Info starts a new message with info level.
	Info = log.Info
	// Warn starts a new message with warn level.
	Warn = log.Warn
	// Error starts a new message with error level.
	Error = log.Error
	// Fatal starts a new message with fatal level. The os.Exit(1) function
	// is called by the Msg method.
	Fatal = log.Fatal
)
