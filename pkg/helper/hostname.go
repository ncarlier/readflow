package helper

import (
	"os"
	"strings"
)

// GetHostname returns hostname
func GetHostname() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}
	return hostname
}

// GetMailHostname return mail hostname
func GetMailHostname() string {
	hostname := GetHostname()
	// ugly convention used in order to avoid another config params
	if strings.HasPrefix(hostname, "api.") {
		hostname = strings.Replace(hostname, "api", "my", 1)
	}
	return hostname
}
