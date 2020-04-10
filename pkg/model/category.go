package model

import (
	"time"
)

// Category structure definition
type Category struct {
	ID        *uint      `json:"id,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Rule      *string    `json:"rule,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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

// Rule set category rule
func (cb *CategoryBuilder) Rule(rule string) *CategoryBuilder {
	cb.category.Rule = &rule
	return cb
}

// UserID set category user ID
func (cb *CategoryBuilder) UserID(userID uint) *CategoryBuilder {
	cb.category.UserID = &userID
	return cb
}
