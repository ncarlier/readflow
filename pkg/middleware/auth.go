package middleware

import (
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
)

const usingAuthNMsg = "using authentication"

// Auth is a middleware to authenticate HTTP request
func Auth(cfg config.AuthNConfig) Middleware {
	logger := log.With().Str("method", cfg.Method).Logger()
	var authn Middleware
	switch cfg.Method {
	case "mock":
		authn = MockAuth
	case "proxy":
		authn = ProxyAuth(cfg)
	case "basic":
		logger = logger.With().Str("htpasswd", cfg.Basic.HtpasswdFile).Logger()
		authn = BasicAuth(cfg)
	case "oidc":
		logger = logger.With().Str("issuer", cfg.OIDC.Issuer).Logger()
		authn = OpenIDConnectJWTAuth(cfg)
	default:
		log.Fatal().Str("method", cfg.Method).Msg("non supported authentication method")
	}
	logger.Info().Msg(usingAuthNMsg)
	return authn
}
