package validator

import (
	"errors"
	"fmt"
	"strings"
)

type FieldsValidator struct {
	errors []string
}

func (v *FieldsValidator) Validate(name string, isValid func() bool) {
	if !isValid() {
		v.errors = append(v.errors, fmt.Sprintf("invalid value for '%s' property", name))
	}
}

func (v *FieldsValidator) Error() error {
	if len(v.errors) > 0 {
		return errors.New(strings.Join(v.errors, "\n"))
	}
	return nil
}
