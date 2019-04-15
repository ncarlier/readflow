package event

import (
	"time"

	"github.com/ncarlier/readflow/pkg/model"
)

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

// Issue event issue
type Issue struct {
	URL  *string   `json:"url,omitempty"`
	Date time.Time `json:"date"`
}

// Event structure definition
type Event struct {
	Action string `json:"action"`
	Issue  Issue  `json:"issue"`
}

// UserEvent is an event for user actions
type UserEvent struct {
	Event
	Payload model.User `json:"payload"`
}
