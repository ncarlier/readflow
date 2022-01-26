package config

import (
	"embed"
	"flag"
	"io"
	"os"
)

// Content assets
//go:embed readflow.toml
var ConfigFile embed.FS

// InitConfigFile initialize the config file
var InitConfigFile = flag.String("init-config", "", "Initialize configuration file")

// WriteConfigFile write configuration file
func WriteConfigFile(filename string) error {
	src, err := ConfigFile.Open("readflow.toml")
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
