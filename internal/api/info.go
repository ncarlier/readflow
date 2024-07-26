package api

import (
	"encoding/json"
	"net/http"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/internal/version"
	"github.com/ncarlier/readflow/pkg/utils"
)

// Info API informations model structure.
type Info struct {
	Version   string `json:"version"`
	Authority string `json:"authority"`
	VAPID     string `json:"vapid"`
}

// info is the handler to show API details.
func info(conf *config.Config) http.Handler {
	v := Info{
		Version: version.Version,
		VAPID:   service.Lookup().GetProperties().VAPIDPublicKey,
	}
	if conf.AuthN.Method == "oidc" {
		v.Authority = conf.AuthN.OIDC.Issuer
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(v)
		if err != nil {
			utils.WriteJSONProblem(w, utils.JSONProblem{
				Detail: err.Error(),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
}
