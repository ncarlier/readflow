package archive

import (
	"context"
	"errors"

	"github.com/ncarlier/reader/pkg/model"
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
	default:
		err = errors.New("archive provider not supported: " + config.Provider)
	}
	return provider, err
}
