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

const (
	maxArticlesToCreatePerRequest = 100
	maxArticlesToScrapePerRequest = 10
)

func assertLimitations(articles []model.ArticleCreateForm) error {
	if len(articles) > maxArticlesToCreatePerRequest {
		return errors.New("too many articles")
	}
	articlesToScrape := 0
	for _, article := range articles {
		if article.URL != nil && !article.IsComplete() {
			articlesToScrape++
		}
	}
	if (articlesToScrape) > maxArticlesToScrapePerRequest {
		return errors.New("too many articles to scrape")
	}
	return nil
}

// articles is the handler to post articles using API keys.
func articles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		articlesForm := []model.ArticleCreateForm{}
		articleForm := model.ArticleCreateForm{}

		// Decode body payload
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
			err = errors.New("unexpected delimiter")
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if simpleObject {
			articlesForm = append(articlesForm, articleForm)
		}

		// Apply limitations
		if err := assertLimitations(articlesForm); err != nil {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}

		// Create articles(s)
		articles := service.Lookup().CreateArticles(ctx, articlesForm)

		// TODO filters some attributes

		// Build response
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
