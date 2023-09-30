package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	jwtRequest "github.com/golang-jwt/jwt/v4/request"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/oidc"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/rs/zerolog/log"
)

// OpenIDConnectJWTAuth is a middleware to checks HTTP request with a valid JWT
func OpenIDConnectJWTAuth(cfg config.AuthNConfig) Middleware {
	admins := strings.Split(cfg.Admins, ",")
	oidcConfig, err := oidc.GetOIDCConfiguration(cfg.OIDC.Issuer)
	if err != nil {
		log.Fatal().Err(err).Str("issuer", cfg.OIDC.Issuer).Msg("unable to get OIDC configuration from authority")
	}
	keystore, err := oidc.NewOIDCKeystore(oidcConfig)
	if err != nil {
		log.Fatal().Err(err).Str("issuer", cfg.OIDC.Issuer).Msg("unable to get OIDC keys from JWKS endpoint")
	}
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			w.Header().Set("WWW-Authenticate", `Bearer realm="readflow"`)

			token, err := jwtRequest.ParseFromRequest(r, jwtRequest.OAuth2Extractor, func(token *jwt.Token) (i interface{}, e error) {
				if id, ok := token.Header["kid"]; ok {
					return keystore.GetKey(id.(string))
				}
				return nil, errors.New("kid header not found in token")
			})
			if err != nil {
				jsonErrors(w, err.Error(), 401)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				username := ""
				if val, ok := claims["preferred_username"]; ok {
					username = val.(string)
				} else if val, ok := claims["email"]; ok {
					username = val.(string)
				}
				if username == "" {
					jsonErrors(w, "No username inside token", 403)
					return
				}
				user, err := service.Lookup().GetOrRegisterUser(ctx, username)
				if err != nil {
					jsonErrors(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ctx = context.WithValue(ctx, constant.ContextUser, *user)
				ctx = context.WithValue(ctx, constant.ContextUserID, *user.ID)
				ctx = context.WithValue(ctx, constant.ContextIsAdmin, helper.ContainsString(admins, username))
				inner.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			jsonErrors(w, "Unauthorized", 401)
		})
	}
}
