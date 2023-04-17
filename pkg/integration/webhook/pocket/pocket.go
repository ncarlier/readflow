package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
)

// pocketEntry is the structure definition of a Pocket entry
type pocketEntry struct {
	Title       string  `json:"title,omitempty"`
	URL         *string `json:"url,omitempty"`
	ConsumerKey string  `json:"consumer_key"`
	AccessToken string  `json:"access_token"`
}

// pocketProviderConfig is the structure definition of a Pocket API configuration
type pocketProviderConfig struct {
	Username string `json:"username"`
}

// pocketProvider is the structure definition of a Pocket webhook provider
type pocketProvider struct {
	config      pocketProviderConfig
	consumerKey string
	accessToken string
}

func newPocketProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := pocketProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate username
	if config.Username == "" {
		return nil, fmt.Errorf("missing username")
	}

	// Validate secrets
	accessToken, ok := srv.Secrets["access_token"]
	if !ok {
		return nil, fmt.Errorf("missing access token")
	}

	provider := &pocketProvider{
		config:      config,
		consumerKey: conf.Integration.Pocket.ConsumerKey,
		accessToken: accessToken,
	}

	return provider, nil
}

// Send article to Pocket endpoint.
func (wp *pocketProvider) Send(ctx context.Context, article model.Article) error {
	entry := pocketEntry{
		Title:       article.Title,
		URL:         article.URL,
		ConsumerKey: wp.consumerKey,
		AccessToken: wp.accessToken,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequest("POST", "https://getpocket.com/v3/add", b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", constant.ContentTypeJSON)
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
	webhook.Register("pocket", &webhook.Def{
		Name:   "Pocket",
		Desc:   "Send article(s) to Pocket instance.",
		Create: newPocketProvider,
	})
}
