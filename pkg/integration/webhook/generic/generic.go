package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
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

// webhookArticle is the structure definition of a generic Webhook article
type webhookArticle struct {
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// ProviderConfig is the structure definition of a Webhook configuration
type ProviderConfig struct {
	Endpoint    string            `json:"endpoint"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	Format      string            `json:"format"`
}

// Provider is the structure definition of a Webhook outbound service
type Provider struct {
	config ProviderConfig
	tpl    *template.Template
}

func newWebhookProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
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
		tplName := fmt.Sprintf("webhook-%s", helper.Hash(config.Format))
		tpl, err = template.New(tplName).Parse(config.Format)
		if err != nil {
			return nil, err
		}
	}

	provider := &Provider{
		config: config,
		tpl:    tpl,
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (whp *Provider) Send(ctx context.Context, article model.Article) error {
	art := webhookArticle{
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

func init() {
	webhook.Register("generic", &webhook.Def{
		Name:   "Generic webhook",
		Desc:   "Export article(s) to a generic webhook.",
		Create: newWebhookProvider,
	})
}
