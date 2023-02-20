package event

const (
	// CreateUser is the create event on an user
	CreateUser = "user:create"
	// UpdateUser is the update event on an user
	UpdateUser = "user:update"
	// DeleteUser is the delete event on an user
	DeleteUser = "user:delete"
	// CreateArticle is the create event on an article
	CreateArticle = "article:create"
	// UpdateArticle is the update event on an article
	UpdateArticle = "article:update"
)

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

const (
	// NoNotification will disable global notification policy
	NoNotification EventOption = 1 << iota // 1
)
