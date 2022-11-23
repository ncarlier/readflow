package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/template"
	_ "github.com/ncarlier/readflow/pkg/template/all"
)

func isValidContentType(contentType string) bool {
	switch true {
	case
		contentType == "",
		strings.HasPrefix(contentType, constant.ContentTypeForm),
		strings.HasPrefix(contentType, "application/json"),
		strings.HasPrefix(contentType, "text/"):
		return true
	}
	return false
}

// ProviderConfig is the structure definition of a Webhook configuration
type ProviderConfig struct {
	Endpoint string            `json:"endpoint"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

// Provider is the structure definition of a Webhook outbound service
type Provider struct {
	config         ProviderConfig
	headers        http.Header
	templateEngine template.Provider
	hrefBase       string
}

func newWebhookProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	_, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Extract headers
	headers := http.Header{}
	for k, v := range config.Headers {
		headers.Set(k, v)
	}

	// Validate Content-Type
	if !isValidContentType(headers.Get("Content-Type")) {
		return nil, fmt.Errorf("Content-Type not supported")
	}

	var templateEngine template.Provider
	if strings.TrimSpace(config.Body) != "" {
		// Use template engine if body is not empty
		templateEngine, err = template.NewTemplateEngine("fast", config.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// force content-type to JSON otherwise
		headers.Set("Content-Type", constant.ContentTypeJSON)
	}

	provider := &Provider{
		config:         config,
		headers:        headers,
		templateEngine: templateEngine,
		hrefBase:       strings.Replace(conf.Global.PublicURL, "api.", "", 1),
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (whp *Provider) Send(ctx context.Context, article model.Article) error {
	data := article.ToMap()
	data["href"] = fmt.Sprintf("%s/inbox/%d", whp.hrefBase, article.ID)

	// Build payload
	b := new(bytes.Buffer)
	if whp.templateEngine != nil {
		if err := whp.templateEngine.Execute(b, data); err != nil {
			return err
		}
	} else {
		if err := json.NewEncoder(b).Encode(data); err != nil {
			return err
		}
	}

	// Build request
	req, err := http.NewRequest("POST", whp.config.Endpoint, b)
	if err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("User-Agent", constant.UserAgent)
	req.Header.Set("Content-Type", constant.ContentTypeJSON)
	// Set configured headers
	for k, v := range whp.headers {
		req.Header.Set(k, v[0])
	}

	// Do HTTP request
	client := &http.Client{Timeout: constant.DefaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return nil
}

func init() {
	webhook.Register("generic", &webhook.Def{
		Name:   "Generic webhook",
		Desc:   "Export article(s) to a generic webhook.",
		Create: newWebhookProvider,
	})
}
