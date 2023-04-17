package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
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
}

// Provider is the structure definition of a Nunux Keeper webhook provider
type Provider struct {
	config   ProviderConfig
	APIKey   string
	endpoint *url.URL
}

func newKeeperProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate secrets
	apiKey, ok := srv.Secrets["api_key"]
	if !ok {
		return nil, fmt.Errorf("missing API key")
	}

	provider := &Provider{
		config:   config,
		APIKey:   apiKey,
		endpoint: endpoint,
	}

	return provider, nil
}

// Send article to Nunux Keeper endpoint.
func (p *Provider) Send(ctx context.Context, article model.Article) error {
	art := Article{
		Title:       article.Title,
		Origin:      article.URL,
		Content:     article.HTML,
		ContentType: "text/html",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(art)

	req, err := http.NewRequest("POST", p.getAPIEndpoint("/v2/documents"), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", constant.ContentTypeJSON)
	req.SetBasicAuth("api", p.APIKey)
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

func (p *Provider) getAPIEndpoint(path string) string {
	baseURL := *p.endpoint
	baseURL.Path = path
	return baseURL.String()
}

func init() {
	webhook.Register("keeper", &webhook.Def{
		Name:   "Nunux Keeper",
		Desc:   "Send article(s) to Nunux Keeper instance.",
		Create: newKeeperProvider,
	})
}
