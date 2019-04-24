package api

import (
	"fmt"
	"net/http"

	"github.com/ncarlier/readflow/pkg/service"

	"github.com/ncarlier/readflow/pkg/config"
)

const vapidTpl = "\nwindow.vapidPublicKey = %q;\n"

func initJS(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "max-age=2592000")
		pubKey := service.Lookup().GetProperties().VAPIDPublicKey
		fmt.Fprintf(w, vapidTpl, pubKey)
	})
}
