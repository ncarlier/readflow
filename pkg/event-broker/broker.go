package eventbroker

import (
	"fmt"
	"io"
	"net/url"

	"github.com/rs/zerolog/log"
)

var instance Broker

// Broker is the event broker interface
type Broker interface {
	Send(payload io.Reader) error
}

// Configure the data event Broker regarding the configuration URI
func Configure(uri string) (Broker, error) {
	if uri == "" {
		return nil, nil
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration URI: %s", uri)
	}

	switch u.Scheme {
	case "http", "https":
		instance, err = newHTTPBroker(u)
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "broker").Str("uri", u.String()).Msg("using HTTP event broker")
	default:
		return nil, fmt.Errorf("unsupported event broker: %s", u.Scheme)
	}
	return instance, nil
}

// Lookup returns the global event broker
func Lookup() Broker {
	if instance != nil {
		return instance
	}
	return nil
}
