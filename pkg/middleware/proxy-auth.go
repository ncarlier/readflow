package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/service"
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

// ProxyAuth is a middleware to checks HTTP request credentials from proxied headers
func ProxyAuth(cfg config.AuthNProxyConfig) Middleware {
	usernameHeaderNames := strings.Split(cfg.Headers, ",")
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
				ctx = context.WithValue(ctx, constant.ContextUser, *user)
				ctx = context.WithValue(ctx, constant.ContextUserID, *user.ID)
				ctx = context.WithValue(ctx, constant.ContextIsAdmin, false)
				inner.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			w.Header().Set("Proxy-Authenticate", `Basic realm="readflow"`)
			jsonErrors(w, "Unauthorized", 401)
		})
	}
}
