package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/oidc"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/rs/zerolog/log"
)

// OpenIDConnectJWTAuth is a middleware to checks HTTP request with a valid JWT
func OpenIDConnectJWTAuth(authority string) Middleware {
	cfg, err := oidc.GetOIDCConfiguration(authority)
	if err != nil {
		log.Fatal().Err(err).Str("authority", authority).Msg("unable to get OIDC configuration from authority")
	}
	keystore := oidc.NewOIDCKeystore(cfg)
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			w.Header().Set("WWW-Authenticate", `Bearer realm="Restricted"`)

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
				isAdmin := false
				if val, ok := claims["realm_access"]; ok {
					if val, ok = val.(map[string]interface{})["roles"]; ok {
						for _, role := range val.([]interface{}) {
							if role.(string) == "admin" {
								isAdmin = true
								break
							}
						}
					}
				}
				ctx = context.WithValue(ctx, constant.ContextUser, *user)
				ctx = context.WithValue(ctx, constant.ContextUserID, *user.ID)
				ctx = context.WithValue(ctx, constant.ContextIsAdmin, isAdmin)
				inner.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			jsonErrors(w, "Unauthorized", 401)
		})
	}
}
