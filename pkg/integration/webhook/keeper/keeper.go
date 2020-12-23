package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
)

// Article is the structure definition of a Nunux Keeper article
type Article struct {
	Title       string  `json:"title,omitempty"`
	Origin      *string `json:"origin,omitempty"`
	Content     *string `json:"content,omitempty"`
	ContentType string  `json:"content_type,omitempty"`
}

// ProviderConfig is the structure definition of a Nunux Keeper API configuration
type ProviderConfig struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
}

// Provider is the structure definition of a Nunux Keeper webhook provider
type Provider struct {
	config ProviderConfig
}

func newKeeperProvider(srv model.OutgoingWebhook) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	_, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		config: config,
	}

	return provider, nil
}

// Send article to Nunux Keeper endpoint.
func (kp *Provider) Send(ctx context.Context, article model.Article) error {
	art := Article{
		Title:       article.Title,
		Origin:      article.URL,
		Content:     article.HTML,
		ContentType: "text/html",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(art)

	req, err := http.NewRequest("POST", kp.config.Endpoint, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("api", kp.config.APIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		}
		return err
	}

	return nil
}

func init() {
	webhook.Register("keeper", &webhook.Def{
		Name:   "Nunux Keeper",
		Desc:   "Send article(s) to Nunux Keeper instance.",
		Create: newKeeperProvider,
	})
}
