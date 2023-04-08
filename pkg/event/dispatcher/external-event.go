package dispatcher

import (
	"bytes"
	"encoding/json"
	"time"
)

// Issue event issue
type Issue struct {
	URL  *string   `json:"url,omitempty"`
	Date time.Time `json:"date"`
}

// ExternalEvent structure definition
type ExternalEvent struct {
	Action  string      `json:"action"`
	Issue   Issue       `json:"issue"`
	Payload interface{} `json:"payload"`
}

// Marshal event to JSON data
func (ev *ExternalEvent) Marshal() *bytes.Buffer {
	result := new(bytes.Buffer)
	json.NewEncoder(result).Encode(ev)
	return result
}

// NewExternalEvent create an external event
func NewExternalEvent(action string, payload interface{}) *ExternalEvent {
	evt := &ExternalEvent{
		Payload: payload,
	}
	evt.Action = action
	evt.Issue = Issue{
		Date: time.Now(),
		// TODO: set proper issuer URL
		// URL:  conf.PublicURL,
	}
	return evt
}
