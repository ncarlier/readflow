package constant

const (
	// Username is the key used to store current username into the request context
	Username = iota
	// UserID is the key used to store current user ID into the request context
	UserID
	// RequestID is the key used to store request ID into the request context
	RequestID
	// APIKeyAlias is the key used to store API key alias into the request context
	APIKeyAlias
	// IsAdmin is true when the user is an admin
	IsAdmin
)
