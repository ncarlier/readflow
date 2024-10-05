package config

import (
	"expvar"
)

var configMap = expvar.NewMap("config")

func exportConfigVar(key, value string) {
	configMap.Set(key, new(expvar.String))
	configMap.Get(key).(*expvar.String).Set(value)
}

// ExportVars export some configuration variables to expvar
func ExportVars(conf *Config) {
	exportConfigVar("authn-method", conf.AuthN.Method)
	exportConfigVar("http-listen-addr", conf.HTTP.ListenAddr)
	exportConfigVar("smtp-listen-addr", conf.SMTP.ListenAddr)
	exportConfigVar("metrics-listen-addr", conf.Metrics.ListenAddr)
	exportConfigVar("http-public-url", conf.HTTP.PublicURL)
	exportConfigVar("ui-public-url", conf.UI.PublicURL)
	exportConfigVar("image-proxy-url", conf.ImageProxy.URL)
}
