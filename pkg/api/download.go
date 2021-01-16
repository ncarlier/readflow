package api

import (
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/service"
)

// download is the handler for downloaging articles.
func download() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/articles/")
		if id == "" {
			http.Error(w, "missing article ID", http.StatusBadRequest)
			return
		}
		idArticle, ok := helper.ConvGQLStringToUint(id)
		if !ok {
			http.Error(w, "invalid article ID", http.StatusBadRequest)
			return
		}

		// Archive the article
		data, err := service.Lookup().ArchiveArticle(r.Context(), idArticle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
