package model

import (
	"time"
)

// Archiver structure definition
type Archiver struct {
	ID        *uint      `json:"id,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Alias     string     `json:"alias,omitempty"`
	IsDefault bool       `json:"is_default,omitempty"`
	Provider  string     `json:"provider,omitempty"`
	Config    string     `json:"config,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// ArchiverCreateForm structure definition
type ArchiverCreateForm struct {
	Alias     string
	IsDefault bool
	Provider  string
	Config    string
}

// ArchiverUpdateForm structure definition
type ArchiverUpdateForm struct {
	ID        uint
	Alias     *string
	IsDefault *bool
	Provider  *string
	Config    *string
}

// ArchiverCreateFormBuilder is a builder to create an Archiver create form
type ArchiverCreateFormBuilder struct {
	form *ArchiverCreateForm
}

// NewArchiverCreateFormBuilder creates new Archiver builder instance
func NewArchiverCreateFormBuilder() ArchiverCreateFormBuilder {
	form := &ArchiverCreateForm{}
	return ArchiverCreateFormBuilder{form}
}

// Build creates the archiver
func (ab *ArchiverCreateFormBuilder) Build() *ArchiverCreateForm {
	return ab.form
}

// Alias set alias
func (ab *ArchiverCreateFormBuilder) Alias(alias string) *ArchiverCreateFormBuilder {
	ab.form.Alias = alias
	return ab
}

// Provider set provider
func (ab *ArchiverCreateFormBuilder) Provider(provider string) *ArchiverCreateFormBuilder {
	ab.form.Provider = provider
	return ab
}

// Config set config
func (ab *ArchiverCreateFormBuilder) Config(config string) *ArchiverCreateFormBuilder {
	ab.form.Config = config
	return ab
}

// IsDefault set is default
func (ab *ArchiverCreateFormBuilder) IsDefault(isDefault bool) *ArchiverCreateFormBuilder {
	ab.form.IsDefault = isDefault
	return ab
}

// Dummy fill archiver with test data
func (ab *ArchiverCreateFormBuilder) Dummy() *ArchiverCreateFormBuilder {
	ab.form.Provider = "dummy"
	ab.form.Config = "{\"foo\": \"bar\"}"
	return ab
}
