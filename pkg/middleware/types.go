package middleware

import "net/http"

type key int

const (
	// ContextRequestID is the key used to store request ID into the request context
	ContextRequestID key = iota
)

// Middleware function definition
type Middleware func(inner http.Handler) http.Handler
