package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/mediatype"
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
	cfg := pocketProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate username
	if cfg.Username == "" {
		return nil, fmt.Errorf("missing username")
	}

	// Validate secrets
	accessToken, ok := srv.Secrets["access_token"]
	if !ok {
		return nil, fmt.Errorf("missing access token")
	}

	provider := &pocketProvider{
		config:      cfg,
		consumerKey: conf.Integration.Pocket.ConsumerKey,
		accessToken: accessToken,
	}

	return provider, nil
}

// Send article to Pocket endpoint.
func (wp *pocketProvider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	entry := pocketEntry{
		Title:       article.Title,
		URL:         article.URL,
		ConsumerKey: wp.consumerKey,
		AccessToken: wp.accessToken,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://getpocket.com/v3/add", b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mediatype.JSON)
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

	link := ""
	if item, ok := obj["item"]; ok {
		itemMap := item.(map[string]interface{})
		link = fmt.Sprintf("https://getpocket.com/read/%s", itemMap["item_id"])
	}
	result := &webhook.Result{
		URL: &link,
	}

	return result, nil
}

func init() {
	webhook.Register("pocket", &webhook.Def{
		Name:   "Pocket",
		Desc:   "Send article(s) to Pocket instance.",
		Create: newPocketProvider,
	})
}
