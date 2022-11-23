package template

import (
	"io"
)

// Template provider interface
type Provider interface {
	Execute(w io.Writer, data map[string]interface{}) error
}
