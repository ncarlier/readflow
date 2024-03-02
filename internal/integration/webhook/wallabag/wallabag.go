package wallabag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/mediatype"
)

// wallabagEntry is the structure definition of a Wallabag article
type wallabagEntry struct {
	Title   string  `json:"title,omitempty"`
	URL     *string `json:"url,omitempty"`
	Content *string `json:"content,omitempty"`
}

// wallabagTokenResponse is the structure of a OAuth token
type wallabagTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

// wallabagProviderConfig is the structure definition of a Wallabag API configuration
type wallabagProviderConfig struct {
	Endpoint string `json:"endpoint"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
}

// wallabagProvider is the structure definition of a Wallabag webhook provider
type wallabagProvider struct {
	config       wallabagProviderConfig
	endpoint     *url.URL
	clientSecret string
	password     string
}

func newWallabagProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	cfg := wallabagProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate secrets
	clientSecret, ok := srv.Secrets["client_secret"]
	if !ok {
		return nil, fmt.Errorf("missing client secret")
	}
	password, ok := srv.Secrets["password"]
	if !ok {
		return nil, fmt.Errorf("missing password")
	}

	// Validate credentials
	if cfg.ClientID == "" || cfg.Username == "" {
		return nil, fmt.Errorf("wallabag: missing credentials")
	}

	provider := &wallabagProvider{
		config:       cfg,
		endpoint:     endpoint,
		clientSecret: clientSecret,
		password:     password,
	}

	return provider, nil
}

// Send article to Wallabag endpoint.
func (wp *wallabagProvider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	token, err := wp.getAccessToken()
	if err != nil {
		return nil, err
	}

	entry := wallabagEntry{
		Title:   article.Title,
		URL:     article.URL,
		Content: article.HTML,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequestWithContext(ctx, "POST", wp.getAPIEndpoint("/api/entries.json"), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mediatype.JSON)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	obj := make(map[string]interface{})
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, nil
	}

	id := uint(obj["id"].(float64))
	link := wp.getAPIEndpoint(fmt.Sprintf("/view/%d", id))
	result := &webhook.Result{
		URL: &link,
	}
	return result, nil
}

func (wp *wallabagProvider) getAPIEndpoint(path string) string {
	baseURL := *wp.endpoint
	baseURL.Path = path
	return baseURL.String()
}

func (wp *wallabagProvider) getAccessToken() (*wallabagTokenResponse, error) {
	values := url.Values{}
	values.Add("grant_type", "password")
	values.Add("client_id", wp.config.ClientID)
	values.Add("client_secret", wp.clientSecret)
	values.Add("username", wp.config.Username)
	values.Add("password", wp.password)

	res, err := http.PostForm(wp.getAPIEndpoint("/oauth/v2/token"), values)
	if err != nil {
		return nil, fmt.Errorf("wallabag: unable to get access token: %v", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("wallabag: request failed, status=%d", res.StatusCode)
	}
	var token wallabagTokenResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&token); err != nil {
		return nil, fmt.Errorf("wallabag: unable to decode token response: %v", err)
	}

	return &token, nil
}

func init() {
	webhook.Register("wallabag", &webhook.Def{
		Name:   "Wallabag",
		Desc:   "Send article(s) to Wallabag instance.",
		Create: newWallabagProvider,
	})
}
