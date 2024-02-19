package config

import (
	"embed"
	"io"
	"os"
)

//go:embed defaults.toml
var configFile embed.FS

// WriteDefaultConfigFile write default configuration file
func (c *Config) WriteDefaultConfigFile(filename string) error {
	src, err := configFile.Open("defaults.toml")
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
