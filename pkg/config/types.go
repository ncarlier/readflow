package config

// Config is the root of the configuration
type Config struct {
	Global       GlobalConfig       `toml:"global"`
	Integration  IntegrationConfig  `toml:"integration"`
	RateLimiting RateLimitingConfig `toml:"rate_limiting"`
	UserPlans    []UserPlan         `toml:"user_plans"`
}

// GlobalConfig is the global configuration section
type GlobalConfig struct {
	AuthN             string `toml:"authn"`
	DatabaseURI       string `toml:"db"`
	ListenAddr        string `toml:"listen_addr"`
	MetricsListenAddr string `toml:"metrics_listen_addr"`
	SMTPListenAddr    string `toml:"smtp_listen_addr"`
	PublicURL         string `toml:"public_url"`
	UILocation        string `toml:"ui"`
	SecretSalt        string `toml:"secret_salt"`
	BlockList         string `toml:"block_list"`
}

// IntegrationConfig is the integration configuration section
type IntegrationConfig struct {
	ExternalEventBrokerURI string `toml:"external_event_broker_uri"`
	ExternalWebScraperURL  string `toml:"external_web_scraper_url"`
	ImageProxyURL          string `toml:"image_proxy_url"`
	AvatarProvider         string `toml:"avatar_provider"`
	SecretsEngineProvider  string `toml:"secrets_engine_provider"`
	Sentry                 SentryConfiguration
	Pocket                 PocketConfiguration
}

// SentryConfiguration is the Sentry's integration configuration
type SentryConfiguration struct {
	DSN string `toml:"dsn_url"`
}

// PocketConfiguration is the Pocket's integration configuration
type PocketConfiguration struct {
	ConsumerKey string `toml:"consumer_key"`
}

// RateLimitingConfig is rate-limiting configuration section
type RateLimitingConfig struct {
	Notification RateLimiting
	Webhook      RateLimiting
}

// RateLimiterConfig is the configuration of a rate-limiter
type RateLimiting struct {
	// Provider of the rate limiting store
	Provider string `toml:"provider"`
	//Tokens allowed per interval
	Tokens int `toml:"tokens"`
	// Interval until tokens reset
	Interval duration `toml:"interval"`
}

// UserPlanConfig is the configuration a a user plan
type UserPlan struct {
	Name                    string   `toml:"name" json:"name"`
	ArticlesLimit           uint     `toml:"articles_limit" json:"articles_limit"`
	CategoriesLimit         uint     `toml:"categories_limit" json:"categories_limit"`
	IncomingWebhooksLimit   uint     `toml:"incoming_webhooks_limit" json:"incoming_webhooks_limit"`
	OutgoingWebhooksLimit   uint     `toml:"outgoing_webhooks_limit" json:"outgoing_webhooks_limit"`
	OutgoingWebhooksTimeout duration `toml:"outgoing_webhooks_timeout" json:"outgoing_webhooks_timeout"`
}

// GetUserPlan return an user plan by its name and fallback to first plan if missing
func (c Config) GetUserPlan(name string) (result *UserPlan) {
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
