package middleware

import (
	"context"
	"net/http"

	"github.com/ncarlier/reader/pkg/constant"
	"github.com/ncarlier/reader/pkg/service"
)

// MockAuth is a middleware to mock HTTP request credentials
func MockAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := service.Lookup().GetOrRegisterUser(ctx, "call@me.neo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(ctx, constant.UserID, *user.ID)
		inner.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
