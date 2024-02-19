package config

import (
	"embed"
	"os"
	"text/template"
)

// Content assets
//
//go:embed ui.js
var UIConfigFile embed.FS

// WriteUIConfigFile write configuration file
func (c *Config) WriteUIConfigFile(filename string) error {
	tmpl, err := template.New("ui.js").ParseFS(UIConfigFile, "ui.js")
	if err != nil {
		return err
	}

	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	return tmpl.Execute(dst, c)
}
