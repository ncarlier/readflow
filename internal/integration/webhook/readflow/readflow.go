package readflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/mediatype"
)

// ProviderConfig is the structure definition of a Readflow API configuration
type ProviderConfig struct {
	Endpoint string `json:"endpoint"`
}

// Provider is the structure definition of a Readflow webhook provider
type Provider struct {
	config   ProviderConfig
	APIKey   string
	endpoint *url.URL
}

func newReadflowProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	cfg := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate secrets
	apiKey, ok := srv.Secrets["api_key"]
	if !ok {
		return nil, fmt.Errorf("missing API key")
	}

	provider := &Provider{
		config:   cfg,
		APIKey:   apiKey,
		endpoint: endpoint,
	}

	return provider, nil
}

// Send article to Readflow endpoint.
func (p *Provider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	builder := model.NewArticleCreateFormBuilder()
	builder.FromArticle(article)
	if value := ctx.Value(global.ContextUser); value != nil {
		user := value.(model.User)
		builder.Origin(user.Username)
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(builder.Build())

	req, err := http.NewRequestWithContext(ctx, "POST", p.getAPIEndpoint("/articles"), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaults.UserAgent)
	req.Header.Set("Content-Type", mediatype.JSON)
	req.SetBasicAuth("api", p.APIKey)
	client := defaults.HTTPClient
	if _, ok := ctx.Deadline(); ok {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		}
		return nil, err
	}

	return &webhook.Result{}, nil
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
