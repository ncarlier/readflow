package config

import (
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	"github.com/ncarlier/readflow/pkg/types"
)

// Config is the root of the configuration
type Config struct {
	Log          LogConfig          `toml:"log"`
	Database     DatabaseConfig     `toml:"database"`
	HTTP         HTTPConfig         `toml:"http"`
	Metrics      MetricsConfig      `toml:"metrics"`
	SMTP         SMTPConfig         `toml:"smtp"`
	AuthN        AuthNConfig        `toml:"authn"`
	UI           UIConfig           `toml:"ui"`
	Hash         HashConfig         `toml:"hash"`
	Downloader   DownloaderConfig   `toml:"downloader"`
	Scraping     ScrapingConfig     `toml:"scraping"`
	Avatar       AvatarConfig       `toml:"avatar"`
	ImageProxy   ImageProxyConfig   `toml:"image_proxy"`
	PDF          PDFConfig          `toml:"pdf"`
	Secrets      SecretsConfig      `toml:"secrets"`
	Event        EventConfig        `toml:"event"`
	Integration  IntegrationConfig  `toml:"integration"`
	RateLimiting RateLimitingConfig `toml:"rate_limiting"`
	UserPlans    []UserPlan         `toml:"user_plans"`
}

// LogConfig for log configuration section
type LogConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// DatabaseConfig for database configuration section
type DatabaseConfig struct {
	URI string `toml:"uri"`
}

// HTTPConfig for HTTP configuration section
type HTTPConfig struct {
	ListenAddr string `toml:"listen_addr"`
	PublicURL  string `toml:"public_url"`
}

// MetricsConfig for metrics configuration section
type MetricsConfig struct {
	ListenAddr string `toml:"listen_addr"`
}

// SMTPConfig for SMTP configuration section
type SMTPConfig struct {
	ListenAddr string `toml:"listen_addr"`
	Hostname   string `toml:"hostname"`
}

// AuthNConfig for authentication configuration section
type AuthNConfig struct {
	Method string `toml:"method"`
	Admins string `toml:"admins"`
	Basic  AuthNBasicConfig
	OIDC   AuthNOIDCConfig
	Proxy  AuthNProxyConfig
}

// AuthNOIDCConfig for OpenID Connect authentication configuration section
type AuthNOIDCConfig struct {
	Issuer       string `toml:"issuer"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

// AuthNProxyConfig for proxy authentication configuration section
type AuthNProxyConfig struct {
	Headers string `toml:"headers"`
}

// AuthBasicConfig for basci authentication configuration section
type AuthNBasicConfig struct {
	HtpasswdFile string `toml:"htpasswd_file"`
}

// UIConfig for UI configuration section
type UIConfig struct {
	Directory string `toml:"directory"`
	PublicURL string `toml:"public_url"`
	ClientID  string `toml:"client_id"`
}

// HashConfig for hash configuration section
type HashConfig struct {
	SecretKey  types.HexString `toml:"secret_key"`
	SecretSalt types.HexString `toml:"secret_salt"`
}

// DownloaderConfig for downloader configuration
type DownloaderConfig struct {
	UserAgent             string         `toml:"user_agent"`
	Timeout               types.Duration `toml:"timeout"`
	Cache                 string         `toml:"cache"`
	MaxConcurentDownloads uint           `toml:"max_concurent_downloads"`
}

// ScrapingForwardProxyConfig for forward proxy configuration
type ScrapingForwardProxyConfig struct {
	Endpoint string
	Hosts    []string
}

// ScrapingConfig for scraping configuration section
type ScrapingConfig struct {
	UserAgent    string                     `toml:"user_agent"`
	Timeout      types.Duration             `toml:"timeout"`
	BlockList    string                     `toml:"block_list"`
	ForwardProxy ScrapingForwardProxyConfig `toml:"forward_proxy"`
}

// AvatarConfig for avatar configuration section
type AvatarConfig struct {
	ServiceProvider string `toml:"service_provider"`
}

// ImageProxyConfig for image proxy configuration section
type ImageProxyConfig struct {
	URL   string `toml:"url"`
	Sizes string `toml:"sizes"`
}

// PDFConfig for PDF configuration section
type PDFConfig struct {
	ServiceProvider string `toml:"service_provider"`
}

// SecretsConfig for secrets configuration section
type SecretsConfig struct {
	ServiceProvider string `toml:"service_provider"`
}

// EventConfig for event configuration section
type EventConfig struct {
	BrokerURI string `toml:"broker_uri"`
}

// IntegrationConfig for integration configuration section
type IntegrationConfig struct {
	Sentry SentryConfiguration
	Pocket PocketConfiguration
}

// SentryConfiguration for Sentry's integration configuration
type SentryConfiguration struct {
	DSN string `toml:"dsn_url"`
}

// PocketConfiguration for Pocket's integration configuration
type PocketConfiguration struct {
	ConsumerKey string `toml:"consumer_key"`
}

// RateLimitingConfig for rate-limiting configuration section
type RateLimitingConfig struct {
	Notification ratelimiter.RateLimiterConfig
	Webhook      ratelimiter.RateLimiterConfig
}

// UserPlanConfig for user-plan configuration sections
type UserPlan struct {
	Name                    string         `toml:"name" json:"name"`
	ArticlesLimit           uint           `toml:"articles_limit" json:"articles_limit"`
	CategoriesLimit         uint           `toml:"categories_limit" json:"categories_limit"`
	IncomingWebhooksLimit   uint           `toml:"incoming_webhooks_limit" json:"incoming_webhooks_limit"`
	OutgoingWebhooksLimit   uint           `toml:"outgoing_webhooks_limit" json:"outgoing_webhooks_limit"`
	OutgoingWebhooksTimeout types.Duration `toml:"outgoing_webhooks_timeout" json:"outgoing_webhooks_timeout"`
	DownloadTimout          types.Duration `toml:"download_timeout" json:"download_timeout"`
}

// GetUserPlan return an user plan by its name and fallback to first plan if missing
func (c *Config) GetUserPlan(name string) (result *UserPlan) {
	if len(c.UserPlans) == 0 {
		return nil
	}
	for _, plan := range c.UserPlans {
		if plan.Name == name {
			return &plan
		}
	}
	// Fallback to first plan
	plan := c.UserPlans[0]
	return &plan
}
