package template

import (
	"fmt"
)

// TemplateEngineCreator function for create a template engine
type TemplateEngineCreator func(text string) (Provider, error)

// Registry of all template engine
var registry = map[string]TemplateEngineCreator{}

// Register a template engine to the registry
func Register(provider string, creator TemplateEngineCreator) {
	registry[provider] = creator
}

// NewTemplateEngine create new template engine
func NewTemplateEngine(providerName, text string) (Provider, error) {
	creator, ok := registry[providerName]
	if !ok {
		return nil, fmt.Errorf("unsupported template engine provider: %s", providerName)
	}
	return creator(text)
}
