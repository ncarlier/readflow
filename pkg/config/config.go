package config

// Config contain global configuration
type Config struct {
	AuthN             string `flag:"authn" desc:"Authentication method (\"mock\", \"proxy\" or OIDC if URL)" default:"https://login.nunux.org/auth/realms/readflow"`
	Broker            string `flag:"broker" desc:"External event broker URI for outgoing events"`
	DB                string `flag:"db" desc:"Database connection string" default:"postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"`
	ImageProxy        string `flag:"image-proxy" desc:"Image proxy service (passthrough if empty)"`
	ListenAddr        string `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	ListenMetricsAddr string `flag:"listen-metrics" desc:"Metrics listen address"`
	LogLevel          string `flag:"log-level" desc:"Log level (debug, info, warn, error)" default:"info"`
	LogPretty         bool   `flag:"log-pretty" desc:"Output human readable logs" default:"false"`
	PocketConsumerKey string `flag:"pocket-consumer-key" desc:"Pocket consumer key"`
	PublicURL         string `flag:"public-url" desc:"Public URL" default:"https://api.readflow.app"`
	SecretSalt        string `flag:"secret-salt" desc:"Secret salt used by hash algorithms" default:"pepper"`
	SentryDSN         string `flag:"sentry-dsn" desc:"Sentry DSN URL"`
	UserPlans         string `flag:"user-plans" desc:"User plans definition file (deactivated if empty)"`
	WebScraping       string `flag:"web-scraping" desc:"Web Scraping service (internal if empty)"`
}
