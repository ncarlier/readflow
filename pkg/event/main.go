package event

import (
	"github.com/asaskevich/EventBus"
)

var bus = EventBus.New()

// Emit trigger an event to events listenners
func Emit(event string, payload ...interface{}) {
	bus.Publish(event, payload...)
}

// Subscribe add a event listenner
func Subscribe(event string, fn interface{}) error {
	return bus.Subscribe(event, fn)
}
