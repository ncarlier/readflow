package dispatcher

import (
	"fmt"
	"net/http"
	"net/url"
)

// HTTPDispatcher structure
type HTTPDispatcher struct {
	uri *url.URL
}

func newHTTPDispatcher(uri *url.URL) (Dispatcher, error) {
	return &HTTPDispatcher{
		uri: uri,
	}, nil
}

// Send the payload to the event broker
func (hb *HTTPDispatcher) Dispatch(event *ExternalEvent) error {
	// TODO add HMAC header signature (X-Broker-Signature)
	resp, err := http.Post(hb.uri.String(), "application/json; charset=utf-8", event.Marshal())
	if err != nil {
		return err
	} else if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}
