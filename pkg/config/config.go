package config

import (
	"flag"
	"os"
	"strconv"
)

// Config contain global configuration
type Config struct {
	ListenAddr *string
	DB         *string
	Broker     *string
	AuthN      *string
	PublicURL  *string
	Version    *bool
	LogPretty  *bool
	LogLevel   *string
	SentryDSN  *string
}

var config = &Config{
	ListenAddr: flag.String("listen", getEnv("LISTEN_ADDR", ":8080"), "HTTP service address"),
	DB:         flag.String("db", getEnv("DB", "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"), "Database connection string"),
	Broker:     flag.String("broker", getEnv("BROKER", ""), "External event broker URI for outgoing events"),
	AuthN:      flag.String("authn", getEnv("AUTHN", "https://login.nunux.org/auth/realms/readflow"), "Authentication method (\"mock\", \"proxy\" or OIDC if URL)"),
	PublicURL:  flag.String("public-url", getEnv("PUBLIC_URL", "https://api.readflow.app"), "Public URL"),
	Version:    flag.Bool("version", false, "Print version"),
	LogPretty:  flag.Bool("log-pretty", getBoolEnv("LOG_PRETTY", false), "Output human readable logs"),
	LogLevel:   flag.String("log-level", getEnv("LOG_LEVEL", "info"), "Log level (debug, info, warn, error)"),
	SentryDSN:  flag.String("sentry-dsn", getEnv("SENTRY_DSN", ""), "Sentry DSN URL"),
}

func init() {
	// set shorthand parameters
	const shorthand = " (shorthand)"
	usage := flag.Lookup("listen").Usage + shorthand
	flag.StringVar(config.ListenAddr, "l", *config.ListenAddr, usage)
	usage = flag.Lookup("version").Usage + shorthand
	flag.BoolVar(config.Version, "v", *config.Version, usage)
}

// Get global configuration
func Get() *Config {
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("APP_" + key); ok {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	strValue := getEnv(key, strconv.Itoa(fallback))
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	strValue := getEnv(key, strconv.FormatBool(fallback))
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return fallback
}
