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
func ExportVars(conf Config) {
	exportConfigVar("addr", conf.ListenAddr)
	exportConfigVar("authn", conf.ListenAddr)
	exportConfigVar("public-url", conf.PublicURL)
	exportConfigVar("image-proxy", conf.ImageProxy)
}
