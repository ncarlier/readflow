package event

// EventHandler is a event handler function
type EventHandler func(event Event)

// Event structure
type Event struct {
	Name    string
	Payload interface{}
	Option  EventOption
}

// NewEvent create new event
func NewEvent(name string, payload interface{}) *Event {
	return &Event{
		Name:    name,
		Payload: payload,
	}
}

// NewEventWithOption create new event with option
func NewEventWithOption(name string, payload interface{}, option EventOption) *Event {
	result := NewEvent(name, payload)
	result.Option = option
	return result
}

// EventOption is a set of event option.
type EventOption byte

// AddIf adds event option if condition is valid
func (opts *EventOption) AddIf(opt EventOption, condition bool) {
	if condition {
		*opts |= opt
	}
}

// Has test if event options contains this option
func (opts EventOption) Has(opt EventOption) bool {
	return opts&opt != 0
}
