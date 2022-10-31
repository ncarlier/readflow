package readflow

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

// ProviderConfig is the structure definition of a Readflow API configuration
type ProviderConfig struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
}

// Provider is the structure definition of a Readflow webhook provider
type Provider struct {
	config   ProviderConfig
	endpoint *url.URL
}

func newReadflowProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate config
	if config.APIKey == "" {
		return nil, fmt.Errorf("missing API key")
	}

	provider := &Provider{
		config:   config,
		endpoint: endpoint,
	}

	return provider, nil
}

// Send article to Readflow endpoint.
func (p *Provider) Send(ctx context.Context, article model.Article) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(article)

	req, err := http.NewRequest("POST", p.getAPIEndpoint("/articles"), b)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req.Header.Set("Content-Type", constant.ContentTypeJSON)
	req.SetBasicAuth("api", p.config.APIKey)
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
	webhook.Register("readflow", &webhook.Def{
		Name:   "Readflow",
		Desc:   "Send article(s) to Readflow instance.",
		Create: newReadflowProvider,
	})
}
