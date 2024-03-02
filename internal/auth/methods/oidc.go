package methods

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	jwtRequest "github.com/golang-jwt/jwt/v4/request"
	"github.com/ncarlier/readflow/internal/auth"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/ncarlier/readflow/pkg/oidc"
	"github.com/ncarlier/readflow/pkg/utils"
)

// newOIDCAuthMiddleware create a middleware to checks HTTP request with a valid OIDC token
func newOIDCAuthMiddleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
	bearerExtractor := &jwtRequest.BearerExtractor{}
	admins := strings.Split(cfg.Admins, ",")
	oidcClient, err := oidc.NewOIDCClient(cfg.OIDC.Issuer, cfg.OIDC.ClientID, cfg.OIDC.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC client form provided issuer: %w", err)
	}
	keyFunc := buildKeyFunc(oidcClient)
	logger.Info().Str("component", "auth").Str("issuer", cfg.OIDC.Issuer).Msg("using OIDC issuer")
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			w.Header().Set("WWW-Authenticate", `Bearer realm="readflow"`)
			tokenString, err := bearerExtractor.ExtractToken(r)
			if err != nil {
				utils.WriteJSONProblem(w, "", err.Error(), http.StatusBadRequest)
				return
			}

			username := ""
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
			switch {
			case err != nil && errors.Is(err, jwt.ErrTokenMalformed):
				// asume that the token is an opaque token
				// validate token using introspection endpoint and try to get username from it
				username, err = getUsernameFormOpaqueToken(tokenString, oidcClient)
			case token != nil && token.Valid:
				// try to get username from JWT
				username, err = getUsernameFromJWT(token)
			default:
				// error or token invalid
				if err == nil {
					err = errors.New("Unauthorized")
				}
			}

			if err != nil {
				utils.WriteJSONProblem(w, "", err.Error(), http.StatusUnauthorized)
				return
			}

			if username == "" {
				// call UserInfo endpoint to retrive username
				// TODO use cache with subject
				username, err = getUsernameFromUserInfo(tokenString, oidcClient)
				if err != nil {
					utils.WriteJSONProblem(w, "", err.Error(), http.StatusForbidden)
					return
				}
			}

			if username == "" {
				utils.WriteJSONProblem(w, "", "unable to retrieve username from OIDC endpoints", http.StatusForbidden)
				return
			}

			user, err := service.Lookup().GetOrRegisterUser(ctx, username)
			if err != nil {
				utils.WriteJSONProblem(w, "", err.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, global.ContextUser, *user)
			ctx = context.WithValue(ctx, global.ContextUserID, *user.ID)
			ctx = context.WithValue(ctx, global.ContextIsAdmin, utils.ContainsString(admins, username))
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

func buildKeyFunc(client *oidc.Client) jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		if id, ok := token.Header["kid"]; ok {
			return client.Keystore.GetKey(id.(string))
		}
		return nil, errors.New("kid header not found in token")
	}
}

func getUsernameFormOpaqueToken(token string, oidcClient *oidc.Client) (username string, err error) {
	introspection, err := oidcClient.Introspect(token)
	if err != nil {
		return
	}
	if !introspection.Active {
		err = errors.New("token is inactive")
		return
	}
	if introspection.PreferredUsername != "" {
		username = introspection.PreferredUsername
	} else {
		username = introspection.Username
	}
	return
}

func getUsernameFromJWT(token *jwt.Token) (username string, err error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if val, ok := claims["preferred_username"]; ok {
			username = val.(string)
		} else if val, ok := claims["email"]; ok {
			username = val.(string)
		}
	}
	return
}

func getUsernameFromUserInfo(token string, oidcClient *oidc.Client) (username string, err error) {
	info, err := oidcClient.UserInfo(token)
	if err != nil {
		return
	}

	if info.PreferredUsername != "" {
		username = info.PreferredUsername
	} else {
		username = info.Email
	}

	return
}

func init() {
	auth.Register("oidc", newOIDCAuthMiddleware)
}
