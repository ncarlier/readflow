package api

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
)

// index is the handler to show API details.
func index(conf *config.Config) http.Handler {
	if conf.Global.UILocation != "" {
		log.Debug().Str("location", conf.Global.UILocation).Msg("serving UI")
		return http.FileServer(http.Dir(conf.Global.UILocation))
	}
	return http.RedirectHandler("/info", http.StatusSeeOther)
}
