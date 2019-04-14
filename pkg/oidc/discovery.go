package oidc

import (
	"encoding/json"
	"net/http"
)

// GetOIDCConfiguration get OIDC configuration from authority discovery endpoint
func GetOIDCConfiguration(authority string) (OIDCConfiguration, error) {
	var cfg = OIDCConfiguration{}

	r1, err := http.Get(authority + "/.well-known/openid-configuration")
	if err != nil {
		return cfg, err
	}
	defer r1.Body.Close()

	err = json.NewDecoder(r1.Body).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	r2, err := http.Get(cfg.JwksURI)
	if err != nil {
		return cfg, err
	}
	defer r2.Body.Close()

	var jwks = JSONWebKeySet{}
	err = json.NewDecoder(r2.Body).Decode(&jwks)
	if err != nil {
		return cfg, err
	}

	cfg.JSONWebKeySet = jwks.Keys

	return cfg, err
}
