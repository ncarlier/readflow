package event

import (
	"github.com/asaskevich/EventBus"
)

var bus = EventBus.New()

// Emit trigger an event to events listenners
func Emit(event string, payload ...interface{}) {
	bus.Publish(event, payload...)
}
