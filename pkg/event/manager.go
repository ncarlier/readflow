package event

import (
	"github.com/asaskevich/EventBus"
)

// Manager is used to handle events pub/sub
type Manager struct {
	bus EventBus.Bus
}

// Publish an event
func (m *Manager) Publish(event *Event) {
	m.bus.Publish(event.Name, event.Payload, event.Option)
}

// Subscribe an event
func (m *Manager) Subscribe(eventName string, fn EventHandler) {
	m.bus.Subscribe(eventName, func(args ...interface{}) {
		if len(args) > 0 {
			payload := args[0]
			var option EventOption
			if len(args) > 1 {
				option, _ = args[1].(EventOption)
			}
			fn(Event{
				Name:    eventName,
				Payload: payload,
				Option:  option,
			})
		}
	})
}

// NewEventManager creta new event manager
func NewEventManager() *Manager {
	return &Manager{
		bus: EventBus.New(),
	}
}
