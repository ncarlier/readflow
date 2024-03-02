package webhook

import (
	"fmt"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/model"
)

// Creator function for create an webhook provider
type Creator func(srv model.OutgoingWebhook, conf config.Config) (Provider, error)

// Def is a webhook provider definition
type Def struct {
	Name   string
	Desc   string
	Create Creator
}

// Registry of all outgoing webhook provider
var Registry = map[string]*Def{}

// Register add webhook provider definition to the registry
func Register(name string, def *Def) {
	Registry[name] = def
}

// NewOutgoingWebhookProvider create new outgoing webhook provider
func NewOutgoingWebhookProvider(webhook model.OutgoingWebhook, conf config.Config) (Provider, error) {
	def, ok := Registry[webhook.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown webhook service provider: %s", webhook.Provider)
	}
	return def.Create(webhook, conf)
}
