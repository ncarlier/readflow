package api

import (
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
)

var downloadProblem = "unable to download article"

// download is the handler for downloading articles.
func download() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/articles/")
		if id == "" {
			utils.WriteJSONProblem(w, downloadProblem, "missing article ID", http.StatusBadRequest)
			return
		}
		idArticle := utils.ConvGraphQLID(id)
		if idArticle == nil {
			utils.WriteJSONProblem(w, downloadProblem, "invalid article ID", http.StatusBadRequest)
			return
		}
		// Extract and validate token parameter
		q := r.URL.Query()
		format := q.Get("f")
		if format == "" {
			format = "html"
		}

		// Archive the article
		asset, err := service.Lookup().DownloadArticle(r.Context(), *idArticle, format)
		if err != nil {
			utils.WriteJSONProblem(w, downloadProblem, err.Error(), http.StatusInternalServerError)
			return
		}

		// Manage response header
		header := http.Header{}
		header.Add("Transfer-Encoding", "chunked")
		// Write response
		w.WriteHeader(http.StatusOK)
		asset.Write(w, header)
	})
}
