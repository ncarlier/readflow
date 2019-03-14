package middleware

import (
	"context"
	"net/http"

	"github.com/ncarlier/reader/pkg/constant"
)

// Auth is a middleware to checks HTTP request credentials
func Auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := r.Header.Get("X-WEBAUTH-USER")
		if user != "" {
			ctx = context.WithValue(ctx, constant.Username, user)
			inner.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Ah ah ah, you didn't say the magic word"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})
}
