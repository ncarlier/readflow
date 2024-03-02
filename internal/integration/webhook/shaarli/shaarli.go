package shaarli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/mediatype"
)

// shaarliEntry is the structure definition of a Shaarli article
type shaarliEntry struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	URL         *string `json:"url,omitempty"`
	Private     bool    `json:"private"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

// shaarliProviderConfig is the structure definition of a Shaarli API configuration
type shaarliProviderConfig struct {
	Endpoint string `json:"endpoint"`
	Private  bool   `json:"private"`
}

// shaarliProvider is the structure definition of a Shaarli webhook provider
type shaarliProvider struct {
	config   shaarliProviderConfig
	endpoint *url.URL
	secret   string
}

func newShaarliProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	cfg := shaarliProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate secrets
	secret, ok := srv.Secrets["secret"]
	if !ok {
		return nil, fmt.Errorf("missing secret")
	}

	provider := &shaarliProvider{
		config:   cfg,
		endpoint: endpoint,
		secret:   secret,
	}

	return provider, nil
}

// Send article to Shaarli endpoint.
func (p *shaarliProvider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	token, err := p.getAccessToken()
	if err != nil {
		return nil, err
	}

	entry := shaarliEntry{
		Title:       article.Title,
		Description: article.Text,
		URL:         article.URL,
		Private:     p.config.Private,
		Created:     article.CreatedAt.Format(time.RFC3339),
		Updated:     article.UpdatedAt.Format(time.RFC3339),
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequestWithContext(ctx, "POST", p.getAPIEndpoint("/api/v1/links"), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaults.UserAgent)
	req.Header.Set("Content-Type", mediatype.JSON)
	req.Header.Set("Authorization", "Bearer "+token)
	client := defaults.HTTPClient
	if _, ok := ctx.Deadline(); ok {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 201 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		}
		return nil, err
	}

	link := p.getAPIEndpoint(resp.Header.Get("Location"))
	result := &webhook.Result{
		URL: &link,
	}
	return result, nil
}

func (p *shaarliProvider) getAPIEndpoint(path string) string {
	baseURL := *p.endpoint
	baseURL.Path = path
	return baseURL.String()
}

func (p *shaarliProvider) getAccessToken() (string, error) {
	claims := new(jwt.StandardClaims)
	claims.IssuedAt = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(p.secret))
}

func init() {
	webhook.Register("shaarli", &webhook.Def{
		Name:   "Shaarli",
		Desc:   "Send article(s) to Shaarli instance.",
		Create: newShaarliProvider,
	})
}
