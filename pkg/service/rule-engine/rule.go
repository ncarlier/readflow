package ruleengine

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/ncarlier/reader/pkg/model"
)

// RuleProcessor define a rule processor
type RuleProcessor struct {
	expression expr.Node
}

// NewRuleProcessor reates a new rule processor
func NewRuleProcessor(rule *model.Rule) (*RuleProcessor, error) {
	p, err := expr.Parse(rule.Rule, expr.Define("article", model.Article{}))
	if err != nil {
		return nil, fmt.Errorf("invalid rule expression: %s", err)
	}
	return &RuleProcessor{
		expression: p,
	}, nil
}
