package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// NewConfig creates a new default configuration
func NewConfig() *Config {
	return &Config{}
}

// LoadFile loads configuration from a TOML file
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
	return nil
}
