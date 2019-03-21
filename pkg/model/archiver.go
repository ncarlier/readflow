package model

import (
	"time"

	"github.com/brianvoe/gofakeit"
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

// ArchiverForm structure definition
type ArchiverForm struct {
	ID        *uint
	Alias     string
	IsDefault bool
	Provider  string
	Config    string
}

// ArchiverBuilder is a builder to create an Archiver
type ArchiverBuilder struct {
	archiver *Archiver
}

// NewArchiverBuilder creates new Archiver builder instance
func NewArchiverBuilder() ArchiverBuilder {
	archiver := &Archiver{}
	return ArchiverBuilder{archiver}
}

// Build creates the archiver
func (ab *ArchiverBuilder) Build() *Archiver {
	return ab.archiver
}

// Random fill archiver with random data
func (ab *ArchiverBuilder) Random() *ArchiverBuilder {
	gofakeit.Seed(0)
	ab.archiver.Alias = gofakeit.Word()
	return ab
}

// UserID set archiver user ID
func (ab *ArchiverBuilder) UserID(userID uint) *ArchiverBuilder {
	ab.archiver.UserID = &userID
	return ab
}

// Form set archiver content using Form object
func (ab *ArchiverBuilder) Form(form *ArchiverForm) *ArchiverBuilder {
	ab.archiver.ID = form.ID
	ab.archiver.Alias = form.Alias
	ab.archiver.Provider = form.Provider
	ab.archiver.Config = form.Config
	ab.archiver.IsDefault = form.IsDefault
	return ab
}
