package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/mediatype"
	"github.com/ncarlier/readflow/pkg/template"
	_ "github.com/ncarlier/readflow/pkg/template/all"
)

func isValidContentType(contentType string) bool {
	switch {
	case
		contentType == "",
		strings.HasPrefix(contentType, mediatype.FormURLEncoded),
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
	cfg := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	_, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	// Extract headers
	headers := http.Header{}
	for k, v := range cfg.Headers {
		headers.Set(k, v)
	}

	// Validate Content-Type
	if !isValidContentType(headers.Get("Content-Type")) {
		return nil, fmt.Errorf("Content-Type not supported")
	}

	var templateEngine template.Provider
	if strings.TrimSpace(cfg.Body) != "" {
		// Use template engine if body is not empty
		templateEngine, err = template.NewTemplateEngine("fast", cfg.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// force content-type to JSON otherwise
		headers.Set("Content-Type", mediatype.JSON)
	}

	provider := &Provider{
		config:         cfg,
		headers:        headers,
		templateEngine: templateEngine,
		hrefBase:       conf.UI.PublicURL,
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (whp *Provider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	data := article.ToMap()
	data["href"] = fmt.Sprintf("%s/inbox/%d", whp.hrefBase, article.ID)

	// Build payload
	b := new(bytes.Buffer)
	if whp.templateEngine != nil {
		if err := whp.templateEngine.Execute(b, data); err != nil {
			return nil, err
		}
	} else {
		if err := json.NewEncoder(b).Encode(data); err != nil {
			return nil, err
		}
	}

	// Build request
	req, err := http.NewRequestWithContext(ctx, "POST", whp.config.Endpoint, b)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("User-Agent", defaults.UserAgent)
	req.Header.Set("Content-Type", mediatype.JSON)
	// Set configured headers
	for k, v := range whp.headers {
		req.Header.Set(k, v[0])
	}

	// Do HTTP request
	client := defaults.HTTPClient
	if _, ok := ctx.Deadline(); ok {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Validate response status
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return buildWebhookResultFromResponse(resp)
}

func buildWebhookResultFromResponse(resp *http.Response) (result *webhook.Result, err error) {
	link := resp.Header.Get("Location")
	result = &webhook.Result{
		URL: &link,
	}
	contentType := resp.Header.Get("Content-type")
	isText := strings.HasPrefix(contentType, "application/json")
	isJSON := strings.HasPrefix(contentType, "text/")
	var body []byte
	switch {
	case isJSON || isText:
		if body, err = io.ReadAll(resp.Body); err != nil {
			return
		}
		fallthrough
	case isText:
		text := string(body)
		result.Text = &text
	case isJSON:
		json.Unmarshal(body, &result.JSON)
	}
	return
}

func init() {
	webhook.Register("generic", &webhook.Def{
		Name:   "Generic webhook",
		Desc:   "Export article(s) to a generic webhook.",
		Create: newWebhookProvider,
	})
}
