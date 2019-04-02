package event

import (
	"github.com/asaskevich/EventBus"
)

var bus = EventBus.New()

const (
	// CreateUser is the create event on an user
	CreateUser = "user:create"
	// UpdateUser is the update event on an user
	UpdateUser = "user:update"
	// DeleteUser is the delete event on an user
	DeleteUser = "user:delete"
	// CreateArticle is the create event on an article
	CreateArticle = "article:create"
)

// Emit trigger an event to events listenners
func Emit(event string, payload ...interface{}) {
	bus.Publish(event, payload...)
}
