package logger

import (
	"encoding/json"
	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/rs/zerolog"
)

type sentryWriter struct {
	dsn string
}

// SentryWriter create Zerolog writer that send error level message to Sentry
func SentryWriter(dsn string) zerolog.LevelWriter {
	raven.SetDSN(dsn)
	return &sentryWriter{dsn: dsn}
}

// Write implements the io.Writer interface.
func (s *sentryWriter) Write(p []byte) (int, error) {
	var event map[string]interface{}
	tags := make(map[string]string)
	err := json.Unmarshal(p, &event)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return 0, err
	}
	err = fmt.Errorf("%v", event[zerolog.ErrorFieldName])
	// Concert event attributes to string
	for key, value := range event {
		if key != zerolog.ErrorFieldName {
			strValue := fmt.Sprintf("%v", value)
			tags[key] = strValue
		}
	}
	raven.CaptureError(err, tags)
	return len(p), nil
}

// WriteLevel implements the LevelWriter interface.
func (s *sentryWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l == zerolog.ErrorLevel || l == zerolog.FatalLevel {
		return s.Write(p)
	}
	return len(p), nil
}
