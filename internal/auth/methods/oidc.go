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

var bearerExtractor = &jwtRequest.BearerExtractor{}

// newOIDCAuthMiddleware create a middleware to checks HTTP request with a valid OIDC token
func newOIDCAuthMiddleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
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
			// retrieve username from access_token
			username, err := getUsernameFromBearer(r, oidcClient, keyFunc)
			if err != nil {
				utils.WriteJSONProblem(w, utils.JSONProblem{
					Detail: err.Error(),
					Status: http.StatusUnauthorized,
					Context: map[string]interface{}{
						"redirect": "/login",
					},
				})
			}

			// retrieve or register user
			user, err := service.Lookup().GetOrRegisterUser(ctx, username)
			if err != nil {
				utils.WriteJSONProblem(w, utils.JSONProblem{
					Detail: err.Error(),
				})
				return
			}

			// build user context
			ctx = context.WithValue(ctx, global.ContextUser, *user)
			ctx = context.WithValue(ctx, global.ContextUserID, *user.ID)
			ctx = context.WithValue(ctx, global.ContextIsAdmin, utils.ContainsString(admins, username))
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

func getUsernameFromBearer(r *http.Request, oidcClient *oidc.Client, keyFunc jwt.Keyfunc) (username string, err error) {
	tokenString, err := bearerExtractor.ExtractToken(r)
	if err != nil {
		return
	}
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
		return
	}

	if username == "" {
		// call UserInfo endpoint to retrive username
		// TODO use cache with subject
		username, err = getUsernameFromUserInfo(tokenString, oidcClient)
		if err != nil {
			return
		}
	}

	if username == "" {
		err = errors.New("unable to retrieve username from OIDC endpoints")
	}

	return
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

func buildKeyFunc(client *oidc.Client) jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		if id, ok := token.Header["kid"]; ok {
			return client.Keystore.GetKey(id.(string))
		}
		return nil, errors.New("kid header not found in token")
	}
}

func init() {
	auth.Register("oidc", newOIDCAuthMiddleware)
}
