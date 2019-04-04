package model

import (
	"time"
)

// Rule structure definition
type Rule struct {
	ID         *uint      `json:"id,omitempty"`
	UserID     *uint      `json:"user_id,omitempty"`
	Alias      string     `json:"alias,omitempty"`
	CategoryID *uint      `json:"category_id,omitempty"`
	Rule       string     `json:"rule,omitempty"`
	Priority   int        `json:"priority,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// RuleForm structure definition
type RuleForm struct {
	ID         *uint
	Alias      string
	CategoryID uint
	Rule       string
	Priority   int
}

// RuleBuilder is a builder to create an Rule
type RuleBuilder struct {
	rule *Rule
}

// NewRuleBuilder creates new Rule builder instance
func NewRuleBuilder() RuleBuilder {
	rule := &Rule{}
	return RuleBuilder{rule}
}

// Build creates the rule
func (rb *RuleBuilder) Build() *Rule {
	return rb.rule
}

// UserID set rule user ID
func (rb *RuleBuilder) UserID(userID uint) *RuleBuilder {
	rb.rule.UserID = &userID
	return rb
}

// Form set rule content using Form object
func (rb *RuleBuilder) Form(form *RuleForm) *RuleBuilder {
	rb.rule.ID = form.ID
	rb.rule.Alias = form.Alias
	rb.rule.CategoryID = &form.CategoryID
	rb.rule.Rule = form.Rule
	rb.rule.Priority = form.Priority
	return rb
}
