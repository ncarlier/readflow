package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/skip2/go-qrcode"

	"github.com/ncarlier/readflow/internal/config"
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
		endpoint := *r.URL
		endpoint.RawPath = "/articles"
		endpoint.RawFragment = ""
		endpoint.RawQuery = ""

		// Build UI outgoing webhook configuration URL
		payload := fmt.Sprintf(
			"%s/settings/integrations/outgoing-webhooks/add?provider=readflow&endpoint=%s&api_key=%s",
			conf.UI.PublicURL,
			url.QueryEscape(endpoint.String()),
			token,
		)

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
