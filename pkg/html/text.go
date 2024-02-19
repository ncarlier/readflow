package html

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var newLineTags = regexp.MustCompile("title|description|meta|h1|h2|h3|h4|h5|h6|p|pre|code|li|dd|dt")

// HTML2Text get HTML as text
func HTML2Text(content string) (string, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(content))
	var text strings.Builder
	token := tokenizer.Token()
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return text.String(), nil
			}
			return "", tokenizer.Err()
		case tt == html.StartTagToken:
			token = tokenizer.Token()
		case tt == html.TextToken:
			if token.Data == "script" {
				continue
			}
			content := html.UnescapeString(string(tokenizer.Text()))
			if content != "" {
				text.WriteString(content)
			}
			if newLineTags.MatchString(token.Data) {
				text.WriteString("\n")
			}
		}
	}
}
