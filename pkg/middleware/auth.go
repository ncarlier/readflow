package middleware

import (
	"strings"

	"github.com/rs/zerolog/log"
)

const usingAuthNMsg = "using authentication"

// Auth is a middleware to authenticate HTTP request
func Auth(method string) Middleware {
	switch {
	case method == "mock":
		log.Info().Str("method", method).Msg(usingAuthNMsg)
		return MockAuth
	case method == "proxy":
		log.Info().Str("method", method).Msg(usingAuthNMsg)
		return ProxyAuth
	case strings.HasPrefix(method, "file://"):
		log.Info().Str("method", "basic").Str("htpasswd", method).Msg(usingAuthNMsg)
		return BasicAuth(method)
	case strings.HasPrefix(method, "https://"):
		log.Info().Str("method", "bearer").Str("authority", method).Msg(usingAuthNMsg)
		return OpenIDConnectJWTAuth(method)
	default:
		log.Fatal().Str("method", method).Msg("non supported authentication method")
		return nil
	}
}
