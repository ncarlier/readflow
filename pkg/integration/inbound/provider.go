package inbound

import (
	"context"
)

// ServiceProvider inbound service provider interface
type ServiceProvider interface{}

// ServiceWithImportSupport is a inbound service with import support
type ServiceWithImportSupport interface {
	Import(ctx context.Context) error
}
