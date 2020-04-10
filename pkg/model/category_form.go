package model

import (
	"github.com/brianvoe/gofakeit"
)

// CategoryCreateForm structure definition
type CategoryCreateForm struct {
	Title string
	Rule  *string
}

// CategoryUpdateForm structure definition
type CategoryUpdateForm struct {
	ID    uint
	Title *string
	Rule  *string
}

// CategoryCreateFormBuilder is a builder to create an CategoryCreateForm
type CategoryCreateFormBuilder struct {
	form *CategoryCreateForm
}

// NewCategoryBuilder creates new Category builder instance
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

// Rule set category rule
func (cb *CategoryCreateFormBuilder) Rule(rule string) *CategoryCreateFormBuilder {
	cb.form.Rule = &rule
	return cb
}
