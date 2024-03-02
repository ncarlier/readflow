package auth

import (
	"fmt"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/middleware"
)

// AuthMiddlewareCreator function for create an Authentication middleware
type AuthMiddlewareCreator func(cfg *config.AuthNConfig) (middleware.Middleware, error)

// Registry of all middlewares
var registry = map[string]AuthMiddlewareCreator{}

// Register an AuthN middleware to the registry
func Register(method string, creator AuthMiddlewareCreator) {
	registry[method] = creator
}

// NewArticleExporter create new article Exporter
func NewAuthMiddleware(cfg *config.AuthNConfig) (middleware.Middleware, error) {
	creator, ok := registry[cfg.Method]
	if !ok {
		return nil, fmt.Errorf("unsupported authentication method: %s", cfg.Method)
	}
	return creator(cfg)
}
