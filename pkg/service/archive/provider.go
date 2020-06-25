package archive

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/model"
)

// Provider archive provider interface
type Provider interface {
	Archive(ctx context.Context, article model.Article) error
}

// NewArchiveProvider creates new archive provider.
func NewArchiveProvider(config model.Archiver) (Provider, error) {
	var provider Provider
	var err error
	switch config.Provider {
	case "keeper":
		provider, err = newKeeperProvider(config)
	case "webhook":
		provider, err = newWebhookProvider(config)
	case "wallabag":
		provider, err = newWallabagProvider(config)
	default:
		err = errors.New("archive provider not supported: " + config.Provider)
	}
	return provider, err
}
