package outbound

import (
	"context"

	"github.com/ncarlier/readflow/pkg/model"
)

// ServiceProvider outbound service provider interface
type ServiceProvider interface {
	Send(ctx context.Context, article model.Article) error
}
