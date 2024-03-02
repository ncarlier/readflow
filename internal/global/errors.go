package global

import (
	"fmt"
)

// InvalidParameterError create invalid parameter error
func InvalidParameterError(name string) error {
	return fmt.Errorf("a parameter is not valid: %s", name)
}

// RequireParameterError create require parameter error
func RequireParameterError(name string) error {
	return fmt.Errorf("a required parameter is missing: %s", name)
}
