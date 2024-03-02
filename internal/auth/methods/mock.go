package methods

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/internal/auth"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
)

// newMockAuthMiddleware create a middleware to mock HTTP request credentials
func newMockAuthMiddleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
	logger.Info().Str("component", "auth").Msg("using Mock Authentification")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := service.Lookup().GetOrRegisterUser(ctx, "call@me.morpheus")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, global.ContextUser, *user)
			ctx = context.WithValue(ctx, global.ContextUserID, *user.ID)
			ctx = context.WithValue(ctx, global.ContextIsAdmin, true)
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

func init() {
	auth.Register("mock", newMockAuthMiddleware)
}
