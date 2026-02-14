package auth

import "net/http"

// BasicAuthTransport is an HTTP transport that adds Basic Authentication headers
type BasicAuthTransport struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return transport.RoundTrip(req)
}
