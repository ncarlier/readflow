package model

import (
	"time"
)

// Properties structure definition
type Properties struct {
	Rev             *uint
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	CreatedAt       *time.Time
}
