package archive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncarlier/readflow/pkg/model"
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

// WallabagProviderConfig is the structure definition of a Wallabag API configuration
type WallabagProviderConfig struct {
	Endpoint     string `json:"endpoint"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// WallabagProvider is the structure definition of a Nunux Wallabag archive provider
type WallabagProvider struct {
	config   WallabagProviderConfig
	endpoint *url.URL
}

func newWallabagProvider(archiver model.Archiver) (Provider, error) {
	config := WallabagProviderConfig{}
	if err := json.Unmarshal([]byte(archiver.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate credentials
	if config.ClientID == "" || config.ClientSecret == "" || config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("wallabag: missing credentials")
	}

	provider := &WallabagProvider{
		config:   config,
		endpoint: endpoint,
	}

	return provider, nil
}

// Archive article to Nunux Wallabag instance.
func (wp *WallabagProvider) Archive(ctx context.Context, article model.Article) error {
	token, err := wp.getAccessToken()
	if err != nil {
		return err
	}

	entry := wallabagEntry{
		Title:   article.Title,
		URL:     article.URL,
		Content: article.HTML,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequest("POST", wp.getAPIEndpoint("/api/entries.json"), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
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

func (wp *WallabagProvider) getAPIEndpoint(path string) string {
	baseURL := *wp.endpoint
	baseURL.Path = path
	return baseURL.String()
}

func (wp *WallabagProvider) getAccessToken() (*wallabagTokenResponse, error) {
	values := url.Values{}
	values.Add("grant_type", "password")
	values.Add("client_id", wp.config.ClientID)
	values.Add("client_secret", wp.config.ClientSecret)
	values.Add("username", wp.config.Username)
	values.Add("password", wp.config.Password)

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
