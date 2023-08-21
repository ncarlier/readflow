package mail

import "github.com/emersion/go-smtp"

// the backend implements SMTP server methods
type Backend struct{}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return NewSession(), nil
}

// NewBackend create new SNTP backend
func NewBackend() *Backend {
	return &Backend{}
}
