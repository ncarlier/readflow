package sanitizer

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func newDefaultMinifier() *minify.M {
	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepEndTags: true,
	})
	return minifier
}
