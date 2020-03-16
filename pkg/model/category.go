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
	Rule      string     `json:"rule,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CategoryForm structure definition
type CategoryForm struct {
	ID    *uint
	Title *string
	Rule  *string
}

// CategoryBuilder is a builder to create an Category
type CategoryBuilder struct {
	category *Category
}

// NewCategoryBuilder creates new Category builder instance
func NewCategoryBuilder() CategoryBuilder {
	category := &Category{}
	return CategoryBuilder{category}
}

// Build creates the category
func (cb *CategoryBuilder) Build() *Category {
	return cb.category
}

// BuildForm creates a category form
func (cb *CategoryBuilder) BuildForm() CategoryForm {
	return CategoryForm{
		ID:    cb.category.ID,
		Rule:  &cb.category.Rule,
		Title: &cb.category.Title,
	}
}

// Random fill category with random data
func (cb *CategoryBuilder) Random() *CategoryBuilder {
	gofakeit.Seed(0)
	cb.category.Title = gofakeit.Word()
	return cb
}

// Rule set category rule
func (cb *CategoryBuilder) Rule(rule string) *CategoryBuilder {
	cb.category.Rule = rule
	return cb
}

// UserID set category user ID
func (cb *CategoryBuilder) UserID(userID uint) *CategoryBuilder {
	cb.category.UserID = &userID
	return cb
}

// Form set category content using Form object
func (cb *CategoryBuilder) Form(form *CategoryForm) *CategoryBuilder {
	cb.category.ID = form.ID
	if form.Title != nil {
		cb.category.Title = *form.Title
	}
	if form.Rule != nil {
		cb.category.Rule = *form.Rule
	}
	return cb
}
