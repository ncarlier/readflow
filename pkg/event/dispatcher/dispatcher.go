package dispatcher

import (
	"fmt"
	"net/url"
)

// Dispatcher is the event dispatcher interface
type Dispatcher interface {
	Dispatch(event *ExternalEvent) error
}

// NewDispatcher create new event dispatcher
func NewDispatcher(uri string) (Dispatcher, error) {
	if uri == "" {
		return nil, nil
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid dispatcher URI: %s", uri)
	}

	switch u.Scheme {
	case "http", "https":
		return newHTTPDispatcher(u)
	default:
		return nil, fmt.Errorf("unsupported event dispatcher: %s", u.Scheme)
	}
}
