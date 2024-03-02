package webhook

import (
	"context"

	"github.com/ncarlier/readflow/internal/model"
)

// Provider outgoing webhook provider interface
type Provider interface {
	Send(ctx context.Context, article model.Article) (*Result, error)
}

// Result of an outgoing webhook call
type Result struct {
	URL  *string                 `json:"url,omitempty"`
	Text *string                 `json:"text,omitempty"`
	JSON *map[string]interface{} `json:"json,omitempty"`
}
