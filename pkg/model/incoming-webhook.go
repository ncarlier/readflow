package model

import (
	"time"

	"github.com/google/uuid"
)

// IncomingWebhook structure definition
type IncomingWebhook struct {
	ID          *uint      `json:"id,omitempty"`
	UserID      uint       `json:"user_id,omitempty"`
	Alias       string     `json:"alias,omitempty"`
	Token       string     `json:"token,omitempty"`
	Script      string     `json:"script,omitempty"`
	LastUsageAt *time.Time `json:"last_usage_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// IncomingWebhookCreateForm structure definition
type IncomingWebhookCreateForm struct {
	Alias  string
	Token  string
	Script string
}

// IncomingWebhookUpdateForm structure definition
type IncomingWebhookUpdateForm struct {
	ID     uint
	Alias  *string
	Script *string
}

// IncomingWebhookCreateFormBuilder is a builder to create an incoming webhook create form
type IncomingWebhookCreateFormBuilder struct {
	form *IncomingWebhookCreateForm
}

// NewIncomingWebhookCreateFormBuilder creates new incoming webhook builder instance
func NewIncomingWebhookCreateFormBuilder() IncomingWebhookCreateFormBuilder {
	form := &IncomingWebhookCreateForm{}
	return IncomingWebhookCreateFormBuilder{form}
}

// Build creates the incoming webhook
func (ab *IncomingWebhookCreateFormBuilder) Build() *IncomingWebhookCreateForm {
	ab.form.Token = uuid.New().String()
	return ab.form
}

// Alias set incoming webhook alias
func (ab *IncomingWebhookCreateFormBuilder) Alias(alias string) *IncomingWebhookCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}

// Script set incoming webhook script
func (ab *IncomingWebhookCreateFormBuilder) Script(script string) *IncomingWebhookCreateFormBuilder {
	ab.form.Script = script
	return ab
}
