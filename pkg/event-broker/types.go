package eventbroker

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"
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

// Buffer get user event buffer
func (ue *UserEvent) Buffer() *bytes.Buffer {
	result := new(bytes.Buffer)
	json.NewEncoder(result).Encode(ue)
	return result
}

// NewUserEvent create a user event
func NewUserEvent(user model.User) *UserEvent {
	evt := &UserEvent{
		Payload: user,
	}
	evt.Action = event.CreateUser
	evt.Issue = Issue{
		Date: time.Now(),
		// TODO: set proper issuer URL
		// URL:  conf.PublicURL,
	}
	return evt
}
