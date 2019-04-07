package api

import (
	"encoding/json"
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

// articles is the handler to post articles using API keys.
func articles(conf *config.Config) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		articlesForm := []model.ArticleForm{}

		if err := json.NewDecoder(r.Body).Decode(&articlesForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		articles, err := service.Lookup().CreateArticles(ctx, articlesForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO filters some attributes

		data, err := json.Marshal(articles)
		if err != nil && data == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if len(articles.Errors) > 0 {
			if len(articles.Articles) > 0 {
				w.WriteHeader(http.StatusPartialContent)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.Write(data)
	})
	return middleware.APIKeyAuth(handler)
}
