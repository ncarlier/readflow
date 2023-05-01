package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

// NewConfig create new configuration
func NewConfig() *Config {
	c := &Config{
		Global: GlobalConfig{
			AuthN:       "https://login.nunux.org/auth/realms/readflow",
			DatabaseURI: "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable",
			ListenAddr:  ":8080",
			PublicURL:   "https://readflow.app",
			SecretSalt:  "pepper",
		},
		Integration: IntegrationConfig{
			AvatarProvider: "https://robohash.org/{seed}?set=set4&size=48x48",
		},
		RateLimiting: RateLimitingConfig{
			Notification: RateLimiting{
				Provider: "none",
			},
			Webhook: RateLimiting{
				Provider: "none",
			},
		},
	}
	return c
}

// LoadFile loads the given config file
func (c *Config) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	blob := os.ExpandEnv(string(data))
	if _, err := toml.Decode(blob, c); err != nil {
		return err
	}

	// Apply default configuration...
	mergo.Merge(c, NewConfig())

	return nil
}
