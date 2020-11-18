package model

import (
	"time"
)

// OutboundService structure definition
type OutboundService struct {
	ID        *uint      `json:"id,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Alias     string     `json:"alias,omitempty"`
	IsDefault bool       `json:"is_default,omitempty"`
	Provider  string     `json:"provider,omitempty"`
	Config    string     `json:"config,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// OutboundServiceCreateForm structure definition
type OutboundServiceCreateForm struct {
	Alias     string
	IsDefault bool
	Provider  string
	Config    string
}

// OutboundServiceUpdateForm structure definition
type OutboundServiceUpdateForm struct {
	ID        uint
	Alias     *string
	IsDefault *bool
	Provider  *string
	Config    *string
}

// OutboundServiceCreateFormBuilder is a builder to create an OutboundService create form
type OutboundServiceCreateFormBuilder struct {
	form *OutboundServiceCreateForm
}

// NewOutboundServiceCreateFormBuilder creates new OutboundService builder instance
func NewOutboundServiceCreateFormBuilder() OutboundServiceCreateFormBuilder {
	form := &OutboundServiceCreateForm{}
	return OutboundServiceCreateFormBuilder{form}
}

// Build creates the outboundService
func (ab *OutboundServiceCreateFormBuilder) Build() *OutboundServiceCreateForm {
	return ab.form
}

// Alias set alias
func (ab *OutboundServiceCreateFormBuilder) Alias(alias string) *OutboundServiceCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}

// Provider set provider
func (ab *OutboundServiceCreateFormBuilder) Provider(provider string) *OutboundServiceCreateFormBuilder {
	ab.form.Provider = provider
	return ab
}

// Config set config
func (ab *OutboundServiceCreateFormBuilder) Config(config string) *OutboundServiceCreateFormBuilder {
	ab.form.Config = config
	return ab
}

// IsDefault set is default
func (ab *OutboundServiceCreateFormBuilder) IsDefault(isDefault bool) *OutboundServiceCreateFormBuilder {
	ab.form.IsDefault = isDefault
	return ab
}

// Dummy fill outbound service with test data
func (ab *OutboundServiceCreateFormBuilder) Dummy() *OutboundServiceCreateFormBuilder {
	ab.form.Provider = "dummy"
	ab.form.Config = "{\"foo\": \"bar\"}"
	return ab
}
