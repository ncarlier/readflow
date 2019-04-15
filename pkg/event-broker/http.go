package eventbroker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPBroker structure
type HTTPBroker struct {
	uri *url.URL
}

func newHTTPBroker(uri *url.URL) (Broker, error) {
	return &HTTPBroker{
		uri: uri,
	}, nil
}

// Send the payload to the event broker
func (hb *HTTPBroker) Send(payload io.Reader) error {
	// TODO add HMAC header signature (X-Broker-Signature)
	resp, err := http.Post(hb.uri.String(), "application/json; charset=utf-8", payload)
	if err != nil {
		return err
	} else if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}
