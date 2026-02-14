package config

// Config is the root configuration for readflow-ctl
type Config struct {
	API   APIConfig   `toml:"api"`
	AuthN AuthNConfig `toml:"authn"`
}

// APIConfig for API connection configuration
type APIConfig struct {
	Endpoint string `toml:"endpoint"`
}

// AuthNConfig for authentication configuration
type AuthNConfig struct {
	Method string           `toml:"method"`
	Basic  AuthNBasicConfig `toml:"basic"`
	OIDC   AuthNOIDCConfig  `toml:"oidc"`
}

// AuthNBasicConfig for Basic authentication configuration
type AuthNBasicConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// AuthNOIDCConfig for OpenID Connect authentication configuration
type AuthNOIDCConfig struct {
	Issuer   string `toml:"issuer"`
	ClientID string `toml:"client_id"`
}
