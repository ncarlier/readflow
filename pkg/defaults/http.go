package defaults

import (
	"net/http"
	"time"
)

// UserAgent used by HTTP client
const UserAgent = "Mozilla/5.0 (compatible; Readflow/1.0; +https://github.com/ncarlier/readflow)"

// Timeout for HTTP requests
const Timeout = time.Duration(5 * time.Second)

// HTTPClient is the default HTTP client
var HTTPClient = &http.Client{
	Timeout: Timeout,
}
