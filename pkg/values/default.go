package values

import (
	"net/url"
	"strconv"
)

// GetIntOrDefault return value as int or the default
func GetIntOrDefault(values url.Values, key string, defaultValue int) int {
	value := values.Get(key)
	if value == "" {
		return defaultValue
	}
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return defaultValue
}
