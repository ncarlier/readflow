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

// KeeperArticle is the structure definition of a Nunux Keeper article
type KeeperArticle struct {
	Title       string  `json:"title,omitempty"`
	Origin      *string `json:"origin,omitempty"`
	Content     *string `json:"content,omitempty"`
	ContentType string  `json:"content_type,omitempty"`
}

// KeeperProviderConfig is the structure definition of a Nunux Keeper API configuration
type KeeperProviderConfig struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
}

// KeeperProvider is the structure definition of a Nunux Keeper archive provider
type KeeperProvider struct {
	config KeeperProviderConfig
}

func newKeeperProvider(archiver model.Archiver) (Provider, error) {
	config := KeeperProviderConfig{}
	if err := json.Unmarshal([]byte(archiver.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	_, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	provider := &KeeperProvider{
		config: config,
	}

	return provider, nil
}

// Archive article to Nunux Keeper instance.
func (kp *KeeperProvider) Archive(ctx context.Context, article model.Article) error {
	art := KeeperArticle{
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
