package constant

type key int

const (
	// ContextUserID is the key used to store current user ID into the request context
	ContextUserID key = iota
	// ContextUser is the key used to store current user into the request context
	ContextUser
	// ContextRequestID is the key used to store request ID into the request context
	ContextRequestID
	// ContextIncomingWebhookAlias is the key used to store incomimg webhook alias into the request context
	ContextIncomingWebhookAlias
	// ContextIsAdmin is true when the user is an admin
	ContextIsAdmin
	// ContextDownloader is the key used to store downloader service (uggly hack)
	ContextDownloader
)
