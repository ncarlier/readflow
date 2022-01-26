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
	exportConfigVar("addr", conf.Global.ListenAddr)
	exportConfigVar("authn", conf.Global.AuthN)
	exportConfigVar("public-url", conf.Global.PublicURL)
	exportConfigVar("image-proxy-url", conf.Integration.ImageProxyURL)
}
