package archive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/model"
)

func isValidContentType(contentType string) bool {
	switch contentType {
	case
		constant.ContentTypeForm,
		constant.ContentTypeHTML,
		constant.ContentTypeJSON,
		constant.ContentTypeText:
		return true
	}
	return false
}

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
	Endpoint    string            `json:"endpoint"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	Format      string            `json:"format"`
}

// WebhookProvider is the structure definition of a Webhook archive provider
type WebhookProvider struct {
	config WebhookProviderConfig
	tpl    *template.Template
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

	// Validate Content-Type
	if !isValidContentType(config.ContentType) {
		config.ContentType = constant.ContentTypeJSON
	}

	// Validate format
	var tpl *template.Template
	if config.Format != "" {
		tplName := fmt.Sprintf("webhook-%d", *archiver.ID)
		tpl, err = template.New(tplName).Parse(config.Format)
		if err != nil {
			return nil, err
		}
	}

	provider := &WebhookProvider{
		config: config,
		tpl:    tpl,
	}

	return provider, nil
}

// Archive article to Webhook endpoint.
func (whp *WebhookProvider) Archive(ctx context.Context, article model.Article) error {
	art := WebhookArticle{
		Title:       article.Title,
		Text:        article.Text,
		HTML:        article.HTML,
		URL:         article.URL,
		Image:       article.Image,
		PublishedAt: article.PublishedAt,
	}

	// Build payload
	b := new(bytes.Buffer)
	if whp.tpl != nil {
		if err := whp.tpl.Execute(b, art); err != nil {
			return err
		}
	} else {
		if err := json.NewEncoder(b).Encode(art); err != nil {
			return err
		}
	}

	// Build request
	req, err := http.NewRequest("POST", whp.config.Endpoint, b)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("User-Agent", constant.UserAgent)
	req.Header.Set("Content-Type", whp.config.ContentType)
	for k, v := range whp.config.Headers {
		req.Header.Set(k, v)
	}

	// Do HTTP request
	client := &http.Client{Timeout: constant.DefaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return nil
}
