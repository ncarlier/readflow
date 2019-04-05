package ruleengine

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/ncarlier/reader/pkg/model"
)

// RuleProcessor define a rule processor
type RuleProcessor struct {
	rule       model.Rule
	expression expr.Node
}

// NewRuleProcessor reates a new rule processor
func NewRuleProcessor(rule model.Rule) (*RuleProcessor, error) {
	p, err := expr.Parse(rule.Rule, expr.Define("article", model.Article{}))
	if err != nil {
		return nil, fmt.Errorf("invalid rule expression: %s", err)
	}
	return &RuleProcessor{
		rule:       rule,
		expression: p,
	}, nil
}

// Apply a rule on an article
func (rp *RuleProcessor) Apply(article *model.Article) (bool, error) {
	env := map[string]interface{}{
		"article": article,
	}
	result, err := rp.expression.Eval(env)
	if err != nil {
		return false, fmt.Errorf("Unable to eval expression on article #%d: %s", article.ID, err)
	}
	match := toBool(result)
	if match {
		article.CategoryID = rp.rule.CategoryID
	}

	return match, nil
}

// ProcessorPipeline is a list of rule processor
type ProcessorPipeline []*RuleProcessor

// NewProcessorsPipeline creates a processor pipeline from a rules set
func NewProcessorsPipeline(rules []model.Rule) (*ProcessorPipeline, error) {
	result := ProcessorPipeline{}
	for _, rule := range rules {
		processor, err := NewRuleProcessor(rule)
		if err != nil {
			return nil, err
		}
		result = append(result, processor)
	}
	return &result, nil
}

// Apply a processor pipeline on an article
func (pp *ProcessorPipeline) Apply(article *model.Article) (bool, error) {
	for _, processor := range *pp {
		applied, err := processor.Apply(article)
		if err != nil {
			return false, err
		}
		if applied {
			return true, nil
		}
	}
	return false, nil
}
