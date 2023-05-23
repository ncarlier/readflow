package middleware

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/rs/zerolog/log"
)

// BasicAuth is a middleware to checks HTTP request credentials from Basic AuthN method
func BasicAuth(location string) Middleware {
	htpasswd, err := helper.NewHtpasswdFromFile(location)
	if err != nil {
		log.Fatal().Err(err).Str("location", location).Msg("unable to read htpasswd file")
	}

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			username, password, ok := r.BasicAuth()
			if ok && htpasswd.Authenticate(username, password) {
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
			w.Header().Set("WWW-Authenticate", `Basic realm="readflow", charset="UTF-8"`)
			jsonErrors(w, "Unauthorized", 401)
		})
	}
}
