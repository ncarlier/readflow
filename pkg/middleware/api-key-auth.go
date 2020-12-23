package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/service"
)

// APIKeyAuth is a middleware to checks HTTP request with API key
func APIKeyAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != "api" || pair[1] == "" {
			http.Error(w, "Not authorized", 401)
			return
		}

		token := pair[1]

		inboundService, err := service.Lookup().GetIncomingWebhookByToken(token)
		if err != nil || inboundService == nil {
			http.Error(w, "Not authorized", 401)
			return
		}
		ctx = context.WithValue(ctx, constant.UserID, inboundService.UserID)
		ctx = context.WithValue(ctx, constant.InboundServiceAlias, inboundService.Alias)

		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
