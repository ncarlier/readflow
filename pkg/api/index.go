package api

import (
	"encoding/json"
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/version"
)

// Info API informations model structure.
type Info struct {
	Version string `json:"version"`
}

// index is the handler to show API details.
func index(conf *config.Config) http.Handler {
	v := Info{
		Version: version.Version,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
}
