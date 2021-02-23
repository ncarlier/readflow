package middleware

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/service"
)

var usernameHeaderKeys = []string{
	"X-WEBAUTH-USER",
	"X-Auth-Username",
	"Remote-User",
	"Remote-Name",
}

func getUsernameFromHeader(header http.Header) string {
	for _, key := range usernameHeaderKeys {
		username := header.Get(key)
		if username != "" {
			return username
		}
	}
	return ""
}

// ProxyAuth is a middleware to checks HTTP request credentials from proxied headers
func ProxyAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := getUsernameFromHeader(r.Header)
		if username != "" {
			user, err := service.Lookup().GetOrRegisterUser(ctx, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, constant.UserID, *user.ID)
			inner.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Ah ah ah, you didn't say the magic word"`)
		jsonErrors(w, "Unauthorized", 401)
	})
}
