package webhook

import (
	"context"

	"github.com/ncarlier/readflow/pkg/model"
)

// Provider outgoing webhook provider interface
type Provider interface {
	Send(ctx context.Context, article model.Article) error
}
