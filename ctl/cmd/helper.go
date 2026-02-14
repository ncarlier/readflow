package cmd

import "strings"

// GetFirstCommand returns the first non-flag argument from the argument list
func GetFirstCommand(args []string) (name string, index int) {
	for idx, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return arg, idx
	}
	return "", -1
}
