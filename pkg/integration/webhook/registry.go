package webhook

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/model"
)

// Creator function for create an webhook provider
type Creator func(srv model.OutgoingWebhook) (Provider, error)

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
func NewOutgoingWebhookProvider(conf model.OutgoingWebhook) (Provider, error) {
	def, ok := Registry[conf.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown webhook service provider: %s", conf.Provider)
	}
	return def.Create(conf)
}
