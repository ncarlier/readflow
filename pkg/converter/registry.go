package converter

import (
	"fmt"
)

// Registry of all converters
var registry = map[string]ArticleConverter{}

// Register an article converter to the registry
func Register(format string, converter ArticleConverter) {
	registry[format] = converter
}

// NewArticleConverter create new article converter
func GetArticleConverter(format string) (ArticleConverter, error) {
	converter, ok := registry[format]
	if !ok {
		return nil, fmt.Errorf("unsupported convertion format: %s", format)
	}
	return converter, nil
}
