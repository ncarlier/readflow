package constant

import "time"

// UserAgent used by HTTP client
const UserAgent = "Mozilla/5.0 (compatible; Readflow/1.0; +https://github.com/ncarlier/readflow)"

// DefaultTimeout for HTTP requests
const DefaultTimeout = time.Duration(5 * time.Second)

const (
	// ContentTypeBinary for binary Content-Type
	ContentTypeBinary = "application/octet-stream"
	// ContentTypeForm for URL encoded form Content-Type
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeJSON for JSON Content-Type
	ContentTypeJSON = "application/json; charset=utf-8"
	// ContentTypeHTML for HTML Content-Type
	ContentTypeHTML = "text/html; charset=utf-8"
	// ContentTypeText for text Content-Type
	ContentTypeText = "text/plain; charset=utf-8"
	// ContentTypeEpub for EPUB Content-Type
	ContentTypeEpub = "application/epub+zip"
)
