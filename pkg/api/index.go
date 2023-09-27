package api

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
)

// index is the handler to show API details.
func index(conf *config.Config) http.Handler {
	if conf.UI.Directory != "" {
		log.Debug().Str("location", conf.UI.Directory).Msg("serving UI")
		return http.FileServer(http.Dir(conf.UI.Directory))
	}
	return http.RedirectHandler("/info", http.StatusSeeOther)
}
