package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	"github.com/ncarlier/readflow/pkg/types"
)

// NewConfig create new configuration
func NewConfig() *Config {
	c := &Config{
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
		Database: DatabaseConfig{
			URI: "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable",
		},
		HTTP: HTTPConfig{
			ListenAddr: ":8080",
			PublicURL:  "http://localhost:8080",
		},
		SMTP: SMTPConfig{
			Hostname: "localhost",
		},
		AuthN: AuthNConfig{
			Method: "mock",
			OIDC: AuthNOIDCConfig{
				Issuer: "https://accounts.readflow.app",
			},
			Basic: AuthNBasicConfig{
				HtpasswdFile: "file://.htpasswd",
			},
			Proxy: AuthNProxyConfig{
				Headers: "X-WEBAUTH-USER, X-Auth-Username, Remote-User, Remote-Name",
			},
		},
		UI: UIConfig{
			PublicURL: "http://localhost:8080",
		},
		Hash: HashConfig{
			SecretKey: types.HexString{
				Value: []byte("secret"),
			},
			SecretSalt: types.HexString{
				Value: []byte("pepper"),
			},
		},
		Avatar: AvatarConfig{
			ServiceProvider: "https://robohash.org/{seed}?set=set4&size=48x48",
		},
		Image: ImageConfig{
			ProxySizes: "320,768",
			Cache:      "boltdb:///tmp/readflow-images.cache?maxSize=256,maxEntries=5000,maxEntrySize=1",
		},
		RateLimiting: RateLimitingConfig{
			Notification: ratelimiter.RateLimiterConfig{
				Provider: "none",
			},
			Webhook: ratelimiter.RateLimiterConfig{
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
