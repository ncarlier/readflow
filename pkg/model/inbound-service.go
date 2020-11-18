package model

import (
	"time"

	"github.com/google/uuid"
)

// InboundService structure definition
type InboundService struct {
	ID          *uint      `json:"id,omitempty"`
	UserID      uint       `json:"user_id,omitempty"`
	Alias       string     `json:"alias,omitempty"`
	Token       string     `json:"token,omitempty"`
	Provider    string     `json:"provider,omitempty"`
	Config      string     `json:"config,omitempty"`
	LastUsageAt *time.Time `json:"last_usage_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// InboundServiceCreateForm structure definition
type InboundServiceCreateForm struct {
	Alias    string
	Token    string
	Provider string
	Config   string
}

// InboundServiceUpdateForm structure definition
type InboundServiceUpdateForm struct {
	ID       uint
	Alias    *string
	Provider *string
	Config   *string
}

// InboundServiceCreateFormBuilder is a builder to create an inbound service create form
type InboundServiceCreateFormBuilder struct {
	form *InboundServiceCreateForm
}

// NewInboundServiceCreateFormBuilder creates new inbound service builder instance
func NewInboundServiceCreateFormBuilder() InboundServiceCreateFormBuilder {
	form := &InboundServiceCreateForm{}
	return InboundServiceCreateFormBuilder{form}
}

// Build creates the inbound service
func (ab *InboundServiceCreateFormBuilder) Build() *InboundServiceCreateForm {
	ab.form.Token = uuid.New().String()
	return ab.form
}

// Alias set inbound service alias
func (ab *InboundServiceCreateFormBuilder) Alias(alias string) *InboundServiceCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}

// Provider set provider
func (ab *InboundServiceCreateFormBuilder) Provider(provider string) *InboundServiceCreateFormBuilder {
	ab.form.Provider = provider
	return ab
}

// Config set config
func (ab *InboundServiceCreateFormBuilder) Config(config string) *InboundServiceCreateFormBuilder {
	ab.form.Config = config
	return ab
}

// Dummy fill inbound service with test data
func (ab *InboundServiceCreateFormBuilder) Dummy() *InboundServiceCreateFormBuilder {
	ab.form.Provider = "dummy"
	ab.form.Config = "{\"foo\": \"bar\"}"
	return ab
}
