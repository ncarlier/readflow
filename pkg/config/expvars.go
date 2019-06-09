package config

import (
	"expvar"
	"time"
)

var conf = expvar.NewMap("config")

func getConfString(val *string) func() interface{} {
	return func() interface{} {
		if val == nil {
			return nil
		}
		return *val
	}
}

func getConfDur(val *time.Duration) func() interface{} {
	return func() interface{} {
		if val == nil {
			return nil
		}
		return val.String()
	}
}

func init() {
	conf.Set("addr", expvar.Func(getConfString(config.ListenAddr)))
	conf.Set("authn", expvar.Func(getConfString(config.AuthN)))
	conf.Set("public-url", expvar.Func(getConfString(config.PublicURL)))
	conf.Set("image-proxy", expvar.Func(getConfString(config.ImageProxy)))
}
