package account

import (
	"fmt"

	"github.com/ncarlier/readflow/internal/config"
)

// Creator function for create an account provider
type Creator func(conf *config.IntegrationConfig) (Provider, error)

// Def is a webhook provider definition
type Def struct {
	Name   string
	Desc   string
	Create Creator
}

// Registry of all account provider
var Registry = map[string]*Def{}

// Register add account provider definition to the registry
func Register(name string, def *Def) {
	Registry[name] = def
}

// NewAccountProvider create new account provider
func NewAccountProvider(name string, conf *config.IntegrationConfig) (Provider, error) {
	def, ok := Registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown account provider: %s", name)
	}
	return def.Create(conf)
}
