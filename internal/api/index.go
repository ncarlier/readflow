package api

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/logger"
)

type SPAHandler struct {
	baseDir string
	handler http.Handler
}

func newSPAHandler(baseDir string) http.Handler {
	return &SPAHandler{
		baseDir: baseDir,
		handler: http.FileServer(http.Dir(baseDir)),
	}
}

func (s SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.baseDir, r.URL.Path)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// serve index.html if path does not exist
		http.ServeFile(w, r, filepath.Join(s.baseDir, "index.html"))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise,serve static files
	s.handler.ServeHTTP(w, r)
}

// index is the handler to show API details.
func index(conf *config.Config) http.Handler {
	if conf.UI.Directory != "" {
		logger.Debug().Str("location", conf.UI.Directory).Msg("serving UI")
		// build UI config file from env variables
		configFilename := path.Join(conf.UI.Directory, "config.js")
		if err := conf.WriteUIConfigFile(configFilename); err != nil {
			logger.Warn().Err(err).Str("filename", configFilename).Msg("unable to generate UI config file")
		}
		return newSPAHandler(conf.UI.Directory)
	}
	return http.RedirectHandler("/info", http.StatusSeeOther)
}
