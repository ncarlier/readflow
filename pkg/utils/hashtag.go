package utils

import "regexp"

var hashtagRe = regexp.MustCompile(`#(\S+)`)

// ExtractHashtags extract hashtags from text
func ExtractHashtags(text string) []string {
	return hashtagRe.FindAllString(text, -1)
}

// ReplaceHashtagsPrefix replace into a string hastags with other prefix
func ReplaceHashtagsPrefix(text, prefix string) string {
	return hashtagRe.ReplaceAllString(text, prefix+"$1")
}
