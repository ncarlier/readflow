package html

import (
	"github.com/microcosm-cc/bluemonday"
)

var policy *bluemonday.Policy

// Sanitize HTML content
func Sanitize(content string) string {
	return policy.Sanitize(content)
}

func init() {
	policy = bluemonday.UGCPolicy()
	policy.AddTargetBlankToFullyQualifiedLinks(true)
	policy.AllowAttrs("width", "height", "src", "allowfullscreen", "sandbox").OnElements("iframe")
	policy.AllowAttrs("srcset", "sizes", "data-src").OnElements("img")
}
