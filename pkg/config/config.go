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
	Version    *bool
	Debug      *bool
	LogDir     *string
}

var config = &Config{
	ListenAddr: flag.String("listen", getEnv("LISTEN_ADDR", ":8080"), "HTTP service address"),
	DB:         flag.String("db", getEnv("DB", "postgres://postgres:testpwd@localhost/reader_test?sslmode=disable"), "Database connection string"),
	Version:    flag.Bool("version", false, "Print version"),
	Debug:      flag.Bool("debug", getBoolEnv("DEBUG", false), "Output debug logs"),
	LogDir:     flag.String("log-dir", getEnv("LOG_DIR", os.TempDir()), "Webhooks execution log directory"),
}

func init() {
	// set shorthand parameters
	const shorthand = " (shorthand)"
	usage := flag.Lookup("listen").Usage + shorthand
	flag.StringVar(config.ListenAddr, "l", *config.ListenAddr, usage)
	usage = flag.Lookup("version").Usage + shorthand
	flag.BoolVar(config.Version, "v", *config.Version, usage)
	usage = flag.Lookup("debug").Usage + shorthand
	flag.BoolVar(config.Debug, "d", *config.Debug, usage)
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
