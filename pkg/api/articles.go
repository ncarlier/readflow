package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

// articles is the handler to post articles using API keys.
func articles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		articlesForm := []model.ArticleCreateForm{}
		articleForm := model.ArticleCreateForm{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		body = bytes.TrimSpace(body)
		simpleObject := true
		switch body[0] {
		case byte('['):
			simpleObject = false
			err = json.Unmarshal(body, &articlesForm)
		case byte('{'):
			err = json.Unmarshal(body, &articleForm)
		default: // ] or }
			err = errors.New("Unexpected delimiter")
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if simpleObject {
			articlesForm = append(articlesForm, articleForm)
		}

		articles := service.Lookup().CreateArticles(ctx, articlesForm)

		// TODO filters some attributes

		data, err := json.Marshal(articles)
		if err != nil && data == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		status := http.StatusNoContent
		if len(articles.Errors) == 0 && len(articles.Articles) > 0 {
			status = http.StatusCreated
		} else if len(articles.Errors) > 0 {
			if len(articles.Articles) > 0 {
				status = http.StatusPartialContent
			} else if len(articles.Errors) == 1 && articles.Errors[0] == model.ErrAlreadyExists {
				status = http.StatusNotModified
			} else {
				status = http.StatusBadRequest
			}
		}
		w.WriteHeader(status)
		w.Write(data)
	})
}
