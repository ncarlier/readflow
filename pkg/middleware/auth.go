package middleware

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

const tpl = "using %s as authentication backend"

// Auth is a middleware to authenticate HTTP request
func Auth(method string) Middleware {
	switch method {
	case "mock":
		log.Info().Msg(fmt.Sprintf(tpl, "Mock"))
		return MockAuth
	case "proxy":
		log.Info().Msg(fmt.Sprintf(tpl, "Proxy"))
		return ProxyAuth
	default:
		log.Info().Str("authority", method).Msg(fmt.Sprintf(tpl, "OpenID Connect"))
		return OpenIDConnectJWTAuth(method)
	}
}
