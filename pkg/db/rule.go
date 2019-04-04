package db

import "github.com/ncarlier/reader/pkg/model"

// RuleRepository is the repository interface to manage Rules
type RuleRepository interface {
	GetRuleByID(id uint) (*model.Rule, error)
	GetRuleByUserIDAndAlias(uid uint, alias string) (*model.Rule, error)
	GetRulesByUserID(uid uint) ([]model.Rule, error)
	CreateOrUpdateRule(rule model.Rule) (*model.Rule, error)
	DeleteRule(rule model.Rule) error
	DeleteRules(uid uint, ids []uint) (int64, error)
}
