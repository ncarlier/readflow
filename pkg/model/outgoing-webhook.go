package model

import (
	"time"
)

// OutgoingWebhook structure definition
type OutgoingWebhook struct {
	ID        *uint      `json:"id,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Alias     string     `json:"alias,omitempty"`
	IsDefault bool       `json:"is_default,omitempty"`
	Provider  string     `json:"provider,omitempty"`
	Config    string     `json:"config,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// OutgoingWebhookCreateForm structure definition
type OutgoingWebhookCreateForm struct {
	Alias     string
	IsDefault bool
	Provider  string
	Config    string
}

// OutgoingWebhookUpdateForm structure definition
type OutgoingWebhookUpdateForm struct {
	ID        uint
	Alias     *string
	IsDefault *bool
	Provider  *string
	Config    *string
}

// OutgoingWebhookCreateFormBuilder is a builder to create an outgoing webhook create form
type OutgoingWebhookCreateFormBuilder struct {
	form *OutgoingWebhookCreateForm
}

// NewOutgoingWebhookCreateFormBuilder creates new outgoing webhook builder instance
func NewOutgoingWebhookCreateFormBuilder() OutgoingWebhookCreateFormBuilder {
	form := &OutgoingWebhookCreateForm{}
	return OutgoingWebhookCreateFormBuilder{form}
}

// Build creates the outgoing webhook create form
func (ab *OutgoingWebhookCreateFormBuilder) Build() *OutgoingWebhookCreateForm {
	return ab.form
}

// Alias set alias
func (ab *OutgoingWebhookCreateFormBuilder) Alias(alias string) *OutgoingWebhookCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}

// Provider set provider
func (ab *OutgoingWebhookCreateFormBuilder) Provider(provider string) *OutgoingWebhookCreateFormBuilder {
	ab.form.Provider = provider
	return ab
}

// Config set config
func (ab *OutgoingWebhookCreateFormBuilder) Config(config string) *OutgoingWebhookCreateFormBuilder {
	ab.form.Config = config
	return ab
}

// IsDefault set is default
func (ab *OutgoingWebhookCreateFormBuilder) IsDefault(isDefault bool) *OutgoingWebhookCreateFormBuilder {
	ab.form.IsDefault = isDefault
	return ab
}

// Dummy fill outgoing webhook with test data
func (ab *OutgoingWebhookCreateFormBuilder) Dummy() *OutgoingWebhookCreateFormBuilder {
	ab.form.Provider = "dummy"
	ab.form.Config = "{\"foo\": \"bar\"}"
	return ab
}
