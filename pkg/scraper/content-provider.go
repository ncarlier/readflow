package scraper

import "context"

// ContentProvider is a content provider interface
type ContentProvider interface {
	Get(ctx context.Context, rawurl string) (*WebPage, error)
	Match(url string) bool
}

// ContentProviders is the registry of all supported content provider
var ContentProviders = map[string]ContentProvider{}

// GetContentProvider return content provider that match the given URL
func GetContentProvider(rawurl string) ContentProvider {
	for _, v := range ContentProviders {
		if v.Match(rawurl) {
			return v
		}
	}
	return nil
}
