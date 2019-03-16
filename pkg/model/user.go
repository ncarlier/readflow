package model

import "time"

// User structure definition
type User struct {
	ID          *uint32    `json:"id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Enabled     bool       `json:"enabled,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}