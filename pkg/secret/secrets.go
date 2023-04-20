package secret

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Secrets store
type Secrets map[string]string

// values returns the JSON-encoded representation
func (s Secrets) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan decodes a JSON-encoded value
func (s *Secrets) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}
	// Unmarshal from json to map[string]string
	*s = make(Secrets)
	if err := json.Unmarshal([]byte(v), s); err != nil {
		return err
	}
	return nil
}
