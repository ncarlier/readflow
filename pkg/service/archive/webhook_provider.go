package archive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ncarlier/readflow/pkg/model"
)

// WebhookArticle is the structure definition of a Webhook article
type WebhookArticle struct {
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// WebhookProviderConfig is the structure definition of a Webhook configuration
type WebhookProviderConfig struct {
	Endpoint string `json:"endpoint"`
}

// WebhookProvider is the structure definition of a Webhook archive provider
type WebhookProvider struct {
	config WebhookProviderConfig
}

func newWebhookProvider(archiver model.Archiver) (Provider, error) {
	config := WebhookProviderConfig{}
	if err := json.Unmarshal([]byte(archiver.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	_, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	provider := &WebhookProvider{
		config: config,
	}

	return provider, nil
}

// Archive article to Webhook endpoint.
func (kp *WebhookProvider) Archive(ctx context.Context, article model.Article) error {
	art := WebhookArticle{
		Title:       article.Title,
		Text:        article.Text,
		HTML:        article.HTML,
		URL:         article.URL,
		Image:       article.Image,
		PublishedAt: article.PublishedAt,
	}
	payload := []WebhookArticle{art}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(payload)

	resp, err := http.Post(kp.config.Endpoint, "application/json", b)
	if err != nil || resp.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		}
		return err
	}

	return nil
}
