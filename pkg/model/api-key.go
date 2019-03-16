package model

import (
	"log"
	"time"

	"github.com/ncarlier/reader/pkg/tooling"
)

// APIKey structure definition
type APIKey struct {
	ID          *uint32    `json:"id,omitempty"`
	UserID      uint32     `json:"user_id,omitempty"`
	Alias       string     `json:"alias,omitempty"`
	Token       string     `json:"token,omitempty"`
	LastUsageAt *time.Time `json:"last_usage_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// APIKeyBuilder is a builder to create an APIKey
type APIKeyBuilder struct {
	apiKey *APIKey
}

// NewAPIKeyBuilder creates new APIKey builder instance
func NewAPIKeyBuilder() APIKeyBuilder {
	apiKey := &APIKey{}
	return APIKeyBuilder{apiKey}
}

// Build creates the apiKey
func (ab *APIKeyBuilder) Build() *APIKey {
	ab.apiKey.Token, _ = tooling.NewUUID()
	log.Println(ab.apiKey)
	return ab.apiKey
}

// UserID set apiKey user ID
func (ab *APIKeyBuilder) UserID(userID uint32) *APIKeyBuilder {
	ab.apiKey.UserID = userID
	return ab
}

// Alias set apiKey alias
func (ab *APIKeyBuilder) Alias(alias string) *APIKeyBuilder {
	ab.apiKey.Alias = alias
	return ab
}
