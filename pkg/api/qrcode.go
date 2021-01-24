package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/skip2/go-qrcode"

	"github.com/ncarlier/readflow/pkg/config"
)

// qrcodeHandler is the handler for generating QR code.
func qrcodeHandler(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// Extract and validate token parameter
		token := q.Get("t")
		if token == "" {
			http.Error(w, "bad parameter", http.StatusBadRequest)
			return
		}

		// Build outgoing webhook endpoint
		u, err := url.Parse(conf.PublicURL)
		if err != nil {
			http.Error(w, "invalid public URL", http.StatusInternalServerError)
			return
		}
		u.Path = "/articles"
		u.User = url.UserPassword("api", token)

		// Build UI outgoing webhook configuration URL
		payload := strings.Replace(conf.PublicURL, "api.", "", 1)
		payload = fmt.Sprintf("%s/settings/integrations/outgoing-webhooks/add?enpoint=%s", payload, url.QueryEscape(u.String()))

		// Build QR code
		png, err := qrcode.Encode(payload, qrcode.Medium, 256)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create proxy response
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(png)
	})
}
