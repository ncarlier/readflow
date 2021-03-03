package api

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/integration/account"

	// import all account providers
	_ "github.com/ncarlier/readflow/pkg/integration/account/all"
)

// linking is the handler used for account linking.
func linking(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stub := strings.TrimPrefix(r.URL.Path, "/linking/")
		if stub == "" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		providerName := path.Dir(stub)
		provider, err := account.NewAccountProvider(providerName, conf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		action := path.Base(stub)
		if action == "request" {
			err = provider.RequestHandler(w, r)
		} else if action == "authorize" {
			err = provider.AuthorizeHandler(w, r)
		} else {
			err = errors.New("action non supported")
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
