package api

import (
	"net/http"
	"path"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
)

// index is the handler to show API details.
func index(conf *config.Config) http.Handler {
	if conf.UI.Directory != "" {
		log.Debug().Str("location", conf.UI.Directory).Msg("serving UI")
		// build UI config file from env variables
		configFilename := path.Join(conf.UI.Directory, "config.js")
		if err := conf.WriteUIConfigFile(configFilename); err != nil {
			log.Fatal().Err(err).Str("filename", configFilename).Msg("failed to create UI config file")
		}
		return http.FileServer(http.Dir(conf.UI.Directory))
	}
	return http.RedirectHandler("/info", http.StatusSeeOther)
}
