package config

// Config contain global configuration
type Config struct {
	ListenAddr        string `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	ListenMetricsAddr string `flag:"listen-metrics" desc:"Metrics listen address"`
	DB                string `flag:"db" desc:"Database connection string" default:"postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"`
	Broker            string `flag:"broker" desc:"External event broker URI for outgoing events"`
	AuthN             string `flag:"authn" desc:"Authentication method (\"mock\", \"proxy\" or OIDC if URL)" default:"https://login.nunux.org/auth/realms/readflow"`
	PublicURL         string `flag:"public-url" desc:"Public URL" default:"https://api.readflow.app"`
	LogPretty         bool   `flag:"log-pretty" desc:"Output human readable logs" default:"false"`
	LogLevel          string `flag:"log-level" desc:"Log level (debug, info, warn, error)" default:"info"`
	SentryDSN         string `flag:"sentry-dsn" desc:"Sentry DSN URL"`
	ImageProxy        string `flag:"image-proxy" desc:"Image proxy service (passthrough if empty)"`
	UserPlans         string `flag:"user-plans" desc:"User plans definition file (deactivated if empty)"`
	WebScraping       string `flag:"web-scraping" desc:"Web Scraping service (internal if empty)"`
}
