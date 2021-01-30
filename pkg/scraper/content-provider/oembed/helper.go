package oembed

import (
	"regexp"
	"strings"
)

var (
	su2re1 = regexp.MustCompile("^(https?://[^/]*?)\\*(.+)$")
	su2re2 = regexp.MustCompile("^(https?://[^/]*?/.*?)\\*(.+)$")
	su2re3 = regexp.MustCompile("^(https?://.*?)\\*$")
	su2re4 = regexp.MustCompile("^http://")
)

// Scheme2Regexp convert expr pattern to regular expression
func Scheme2Regexp(scheme string) (*regexp.Regexp, error) {
	expr := strings.Replace(scheme, "?", "\\?", -1)
	expr = su2re1.ReplaceAllString(expr, "${1}[^/]%?${2}")
	expr = su2re2.ReplaceAllString(expr, "${1}.%?${2}")
	expr = su2re3.ReplaceAllString(expr, "${1}.%")
	expr = su2re4.ReplaceAllString(expr, "https?://")
	expr = strings.Replace(expr, "%", "*", -1)
	return regexp.Compile("^" + expr + "$")
}
