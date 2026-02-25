package auth

import (
	"fmt"
	"net/http"
)

// Config holds the authentication parameters resolved at startup.
type Config struct {
	// Method is the deduced auth method: "basic" or "oidc".
	Method string
	// Issuer is the OIDC issuer URL (only for method "oidc").
	Issuer string
	// Username is the Basic auth username (only for method "basic").
	Username string
	// Password is the Basic auth password (only for method "basic").
	Password string
	// ClientID is the OIDC client ID (only for method "oidc").
	ClientID string
}

// NewHTTPClient creates an HTTP client with the appropriate authentication transport.
func NewHTTPClient(conf *Config) (*http.Client, error) {
	switch conf.Method {
	case "basic":
		return &http.Client{
			Transport: &BasicAuthTransport{
				Username: conf.Username,
				Password: conf.Password,
			},
		}, nil
	case "oidc":
		return &http.Client{
			Transport: &OIDCTransport{
				Issuer:   conf.Issuer,
				ClientID: conf.ClientID,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported authentication method: %q", conf.Method)
	}
}
