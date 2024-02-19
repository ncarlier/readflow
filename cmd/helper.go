package cmd

import "strings"

// GetFirstCommand restun first command of argument list
func GetFirstCommand(args []string) (name string, index int) {
	for idx, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// ignore flags
			continue
		}
		return arg, idx
	}
	return "", -1
}
