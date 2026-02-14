package auth

import (
	"fmt"
	"net/http"

	"github.com/ncarlier/readflow/ctl/config"
)

// NewHTTPClient creates an HTTP client with the appropriate authentication transport
func NewHTTPClient(conf *config.Config) (*http.Client, error) {
	switch conf.AuthN.Method {
	case "basic":
		return &http.Client{
			Transport: &BasicAuthTransport{
				Username: conf.AuthN.Basic.Username,
				Password: conf.AuthN.Basic.Password,
			},
		}, nil
	case "oidc":
		return &http.Client{
			Transport: &OIDCTransport{
				Issuer:   conf.AuthN.OIDC.Issuer,
				ClientID: conf.AuthN.OIDC.ClientID,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported authentication method: %q", conf.AuthN.Method)
	}
}
