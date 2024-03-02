package model

import (
	"time"

	"github.com/ncarlier/readflow/pkg/secret"
)

// OutgoingWebhook structure definition
type OutgoingWebhook struct {
	ID        *uint          `json:"id,omitempty"`
	UserID    *uint          `json:"user_id,omitempty"`
	Alias     string         `json:"alias,omitempty"`
	IsDefault bool           `json:"is_default,omitempty"`
	Provider  string         `json:"provider,omitempty"`
	Config    string         `json:"config,omitempty"`
	Secrets   secret.Secrets `json:"-"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
}

// OutgoingWebhookCreateForm structure definition
type OutgoingWebhookCreateForm struct {
	Alias     string
	IsDefault bool
	Provider  string
	Config    string
	Secrets   secret.Secrets
}

// OutgoingWebhookUpdateForm structure definition
type OutgoingWebhookUpdateForm struct {
	ID        uint
	Alias     *string
	IsDefault *bool
	Provider  *string
	Config    *string
	Secrets   *secret.Secrets
}

// OutgoingWebhookCreateFormBuilder is a builder to create an outgoing webhook create form
type OutgoingWebhookCreateFormBuilder struct {
	form *OutgoingWebhookCreateForm
}

// NewOutgoingWebhookCreateFormBuilder creates new outgoing webhook builder instance
func NewOutgoingWebhookCreateFormBuilder() OutgoingWebhookCreateFormBuilder {
	form := &OutgoingWebhookCreateForm{
		Secrets: make(secret.Secrets),
	}
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

// Secrets set secrets
func (ab *OutgoingWebhookCreateFormBuilder) Secrets(secrets secret.Secrets) *OutgoingWebhookCreateFormBuilder {
	ab.form.Secrets = secrets
	return ab
}

// IsDefault set is default
func (ab *OutgoingWebhookCreateFormBuilder) IsDefault(isDefault bool) *OutgoingWebhookCreateFormBuilder {
	ab.form.IsDefault = isDefault
	return ab
}

// Dummy fill outgoing webhook with test data
func (ab *OutgoingWebhookCreateFormBuilder) Dummy() *OutgoingWebhookCreateFormBuilder {
	ab.form.Provider = "generic"
	ab.form.Config = "{\"endpoint\": \"http://example.org\"}"
	ab.form.Secrets["foo"] = "bar"
	return ab
}
