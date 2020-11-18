package inbound

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/model"
)

// ServiceType type of inbound service
type ServiceType string

const (
	// Push inbound service (aka a service that push articles on the webhook)
	Push ServiceType = "push"
	// Pull inbound service (aka a service used to import articles)
	Pull ServiceType = "pull"
)

// ServiceCreator function for create an inbound service provider
type ServiceCreator func(srv model.InboundService) (ServiceProvider, error)

// Service is a inbound service
type Service struct {
	Name   string
	Desc   string
	Type   ServiceType
	Create ServiceCreator
}

// Services registry
var Services = map[string]*Service{}

// Add output to the registry
func Add(name string, service *Service) {
	Services[name] = service
}

// NewInboundServiceProvider create new inbound service provider
func NewInboundServiceProvider(srv model.InboundService) (ServiceProvider, error) {
	service, ok := Services[srv.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown inbound service provider: %s", srv.Provider)
	}
	return service.Create(srv)
}
