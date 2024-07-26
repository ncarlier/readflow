package methods

import (
	"context"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/internal/auth"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/ncarlier/readflow/pkg/utils"
)

func getUsernameFromHeader(header http.Header, keys []string) string {
	for _, key := range keys {
		username := header.Get(strings.TrimSpace(key))
		if username != "" {
			return username
		}
	}
	return ""
}

// newProxyAuthMiddleware create a middleware to checks HTTP request credentials from proxied headers
func newProxyAuthMiddleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
	admins := strings.Split(cfg.Admins, ",")
	usernameHeaderNames := strings.Split(cfg.Proxy.Headers, ",")
	logger.Info().Str("component", "auth").Str("headers", cfg.Proxy.Headers).Msg("using Proxy authentification")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			username := getUsernameFromHeader(r.Header, usernameHeaderNames)
			if username != "" {
				user, err := service.Lookup().GetOrRegisterUser(ctx, username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				ctx = context.WithValue(ctx, global.ContextUser, *user)
				ctx = context.WithValue(ctx, global.ContextUserID, *user.ID)
				ctx = context.WithValue(ctx, global.ContextIsAdmin, utils.ContainsString(admins, username))
				inner.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			w.Header().Set("Proxy-Authenticate", `Basic realm="readflow"`)
			utils.WriteJSONProblem(w, utils.JSONProblem{
				Detail: "invalid authentication header",
				Status: http.StatusUnauthorized,
			})
		})
	}, nil
}

func init() {
	auth.Register("proxy", newProxyAuthMiddleware)
}
