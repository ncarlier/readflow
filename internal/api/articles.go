package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
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

func getStatusCodeFromCreatedArticlesResponse(resp *model.CreatedArticlesResponse) int {
	if len(resp.Articles) > 0 {
		if len(resp.Errors) > 0 {
			// some articles created, some not
			return http.StatusPartialContent
		} else {
			// article(s) created
			return http.StatusCreated
		}
	}
	var globalError error
	for _, err := range resp.Errors {
		if globalError != nil && globalError != err {
			// send internal error if errors are different
			return http.StatusInternalServerError
		}
		globalError = err
	}
	switch globalError {
	case nil:
		// no errors
		return http.StatusNoContent
	case model.ErrAlreadyExists:
		// article already exists
		return http.StatusNotModified
	case service.ErrUserQuotaReached:
		// quota reached
		return http.StatusPaymentRequired
	default:
		// unknown error
		return http.StatusInternalServerError
	}
}

// articles is the handler to post articles using API keys.
func articles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		articlesForm := []model.ArticleCreateForm{}
		articleForm := model.ArticleCreateForm{}

		// Decode body payload
		body, err := io.ReadAll(r.Body)
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
		w.WriteHeader(getStatusCodeFromCreatedArticlesResponse(articles))
		w.Write(data)
	})
}
