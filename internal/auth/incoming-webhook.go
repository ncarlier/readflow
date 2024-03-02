package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/service"
)

// IncomingWebhookAuth is a middleware to checks HTTP request with API key of an incoming webhook
func IncomingWebhookAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if pair[0] != "api" || pair[1] == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		token := pair[1]

		incomingWebhook, err := service.Lookup().GetIncomingWebhookByToken(token)
		if err != nil || incomingWebhook == nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, global.ContextUserID, incomingWebhook.UserID)
		ctx = context.WithValue(ctx, global.ContextIncomingWebhook, incomingWebhook)
		ctx = context.WithValue(ctx, global.ContextIsAdmin, false)

		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
