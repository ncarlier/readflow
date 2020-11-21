package constant

const (
	// Username is the key used to store current username into the request context
	Username = iota
	// UserID is the key used to store current user ID into the request context
	UserID
	// RequestID is the key used to store request ID into the request context
	RequestID
	// InboundServiceAlias is the key used to store inbound service alias into the request context
	InboundServiceAlias
	// IsAdmin is true when the user is an admin
	IsAdmin
)
