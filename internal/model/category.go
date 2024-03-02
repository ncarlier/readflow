package model

import (
	"time"

	"github.com/brianvoe/gofakeit"
)

// Category structure definition
type Category struct {
	ID        *uint      `json:"id,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CategoryCreateForm structure definition
type CategoryCreateForm struct {
	Title string
}

// CategoryUpdateForm structure definition
type CategoryUpdateForm struct {
	ID    uint
	Title *string
}

// CategoryCreateFormBuilder is a builder to create an CategoryCreateForm
type CategoryCreateFormBuilder struct {
	form *CategoryCreateForm
}

// NewCategoryCreateFormBuilder creates new Category create form builder instance
func NewCategoryCreateFormBuilder() CategoryCreateFormBuilder {
	form := &CategoryCreateForm{}
	return CategoryCreateFormBuilder{form}
}

// Build creates the category create form
func (cb *CategoryCreateFormBuilder) Build() *CategoryCreateForm {
	return cb.form
}

// Random fill category with random data
func (cb *CategoryCreateFormBuilder) Random() *CategoryCreateFormBuilder {
	gofakeit.Seed(0)
	cb.form.Title = gofakeit.Word()
	return cb
}
