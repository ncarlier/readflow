package oidc

import (
	"encoding/json"
	"net/http"
)

// GetOIDCConfiguration get OIDC configuration from issuer discovery endpoint
func GetOIDCConfiguration(issuer string) (Configuration, error) {
	var cfg = Configuration{}

	resp, err := http.Get(issuer + "/.well-known/openid-configuration")
	if err != nil {
		return cfg, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
