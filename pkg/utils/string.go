package utils

import (
	"bytes"
	"strings"
)

// OneIsEmpty test if one of the pointers is nil or reference an empty string
func OneIsEmpty(values ...*string) bool {
	for _, value := range values {
		if value == nil || strings.TrimSpace(*value) != "" {
			return true
		}
	}
	return false
}

// Truncate string
func Truncate(value string, size int) string {
	runes := bytes.Runes([]byte(value))
	if len(runes) > size {
		return string(runes[:size]) + "..."
	}
	return string(runes)
}

// IsNilOrEmpty test if string ptr is null or empty
func IsNilOrEmpty(value *string) bool {
	return value == nil || *value == ""
}
