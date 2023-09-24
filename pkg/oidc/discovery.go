package oidc

import (
	"encoding/json"
	"net/http"
)

// GetOIDCConfiguration get OIDC configuration from authority discovery endpoint
func GetOIDCConfiguration(authority string) (Configuration, error) {
	var cfg = Configuration{}

	resp, err := http.Get(authority + "/.well-known/openid-configuration")
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
