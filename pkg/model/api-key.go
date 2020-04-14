package model

import (
	"time"

	"github.com/google/uuid"
)

// APIKeyCreateForm structure definition
type APIKeyCreateForm struct {
	Alias string
	Token string
}

// APIKeyUpdateForm structure definition
type APIKeyUpdateForm struct {
	ID    uint
	Alias string
}

// APIKey structure definition
type APIKey struct {
	ID          *uint      `json:"id,omitempty"`
	UserID      uint       `json:"user_id,omitempty"`
	Alias       string     `json:"alias,omitempty"`
	Token       string     `json:"token,omitempty"`
	LastUsageAt *time.Time `json:"last_usage_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// APIKeyCreateFormBuilder is a builder to create an APIKey create form
type APIKeyCreateFormBuilder struct {
	form *APIKeyCreateForm
}

// NewAPIKeyCreateFormBuilder creates new APIKey builder instance
func NewAPIKeyCreateFormBuilder() APIKeyCreateFormBuilder {
	form := &APIKeyCreateForm{}
	return APIKeyCreateFormBuilder{form}
}

// Build creates the apiKey
func (ab *APIKeyCreateFormBuilder) Build() *APIKeyCreateForm {
	ab.form.Token = uuid.New().String()
	return ab.form
}

// Alias set apiKey alias
func (ab *APIKeyCreateFormBuilder) Alias(alias string) *APIKeyCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}
