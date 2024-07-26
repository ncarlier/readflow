package methods

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/internal/auth"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/htpasswd"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/ncarlier/readflow/pkg/utils"
)

// newBasicAuthMiddlleware create a middleware to checks HTTP request credentials from Basic AuthN method
func newBasicAuthMiddlleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
	admins := strings.Split(cfg.Admins, ",")
	credentials, err := htpasswd.NewHtpasswdFromFile(cfg.Basic.HtpasswdFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read htpasswd file: %w", err)
	}

	logger.Info().Str("component", "auth").Str("htpasswd", cfg.Basic.HtpasswdFile).Msg("using Basic Authentification")

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			username, password, ok := r.BasicAuth()
			if ok && credentials.Authenticate(username, password) {
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
			w.Header().Set("WWW-Authenticate", `Basic realm="readflow", charset="UTF-8"`)
			utils.WriteJSONProblem(w, utils.JSONProblem{
				Detail: "invalid credentials",
				Status: http.StatusUnauthorized,
			})
		})
	}, nil
}

func init() {
	auth.Register("basic", newBasicAuthMiddlleware)
}
