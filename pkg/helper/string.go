package helper

import "strings"

// OneIsEmpty test if one of the pointers is nil or reference an empty string
func OneIsEmpty(values ...*string) bool {
	for _, value := range values {
		if value == nil || len(strings.TrimSpace(*value)) != 0 {
			return true
		}
	}
	return false
}
