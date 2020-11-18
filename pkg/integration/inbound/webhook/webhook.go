package webhook

import (
	"encoding/json"

	"github.com/ncarlier/readflow/pkg/integration/inbound"
	"github.com/ncarlier/readflow/pkg/model"
)

// webhookProviderConfig is the structure definition of a Webhook configuration
type webhookProviderConfig struct{}

// webhookProvider is the structure definition of a Webhook archive provider
type webhookProvider struct {
	config webhookProviderConfig
}

func newWebhookProvider(srv model.InboundService) (inbound.ServiceProvider, error) {
	config := webhookProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}
	provider := &webhookProvider{
		config: config,
	}

	return provider, nil
}

func init() {
	inbound.Add("webhook", &inbound.Service{
		Name:   "Webhook",
		Desc:   "Export article(s) to a webhook.",
		Type:   "Push",
		Create: newWebhookProvider,
	})
}
