package outbound

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/model"
)

// ServiceCreator function for create an outbound service provider
type ServiceCreator func(srv model.OutboundService) (ServiceProvider, error)

// Service is a outbound service
type Service struct {
	Name   string
	Desc   string
	Create ServiceCreator
}

// Services registry
var Services = map[string]*Service{}

// Add output to the registry
func Add(name string, service *Service) {
	Services[name] = service
}

// NewOutboundServiceProvider create new outbound service provider
func NewOutboundServiceProvider(srv model.OutboundService) (ServiceProvider, error) {
	service, ok := Services[srv.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown outbound service provider: %s", srv.Provider)
	}
	return service.Create(srv)
}
