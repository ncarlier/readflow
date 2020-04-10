package ruleengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/constant"

	"github.com/antonmedv/expr"
	"github.com/ncarlier/readflow/pkg/model"
)

// RuleProcessor define a rule processor
type RuleProcessor struct {
	category   model.Category
	expression expr.Node
}

// NewRuleProcessor creates a new rule processor
func NewRuleProcessor(category model.Category) (*RuleProcessor, error) {
	if category.Rule == nil || *category.Rule == "" {
		return nil, nil
	}
	p, err := expr.Parse(
		*category.Rule,
		expr.Define("title", ""),
		expr.Define("text", ""),
		expr.Define("url", ""),
		expr.Define("key", ""),
		expr.Define("tags", []string{}),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid rule expression: %s", err)
	}
	return &RuleProcessor{
		category:   category,
		expression: p,
	}, nil
}

// Apply a rule on an article
func (rp *RuleProcessor) Apply(ctx context.Context, article *model.ArticleForm) (bool, error) {
	tags := []string{}
	if article.Tags != nil {
		tags = strings.Split(*article.Tags, ",")
	}
	text := ""
	if article.Text != nil {
		text = *article.Text
	}
	url := ""
	if article.URL != nil {
		url = *article.URL
	}
	key := ""
	if alias := ctx.Value(constant.APIKeyAlias); alias != nil {
		key = alias.(string)
	}

	env := map[string]interface{}{
		"title": article.Title,
		"text":  text,
		"url":   url,
		"key":   key,
		"tags":  tags,
	}

	result, err := rp.expression.Eval(env)
	if err != nil {
		return false, fmt.Errorf("Unable to eval expression on article %s: %s", *article.URL, err)
	}
	match := toBool(result)
	if match {
		article.CategoryID = rp.category.ID
	}

	return match, nil
}

// ProcessorPipeline is a list of rule processor
type ProcessorPipeline []*RuleProcessor

// NewProcessorsPipeline creates a processor pipeline from a category set
func NewProcessorsPipeline(categories []model.Category) (*ProcessorPipeline, error) {
	result := ProcessorPipeline{}
	for _, category := range categories {
		processor, err := NewRuleProcessor(category)
		if err != nil {
			return nil, err
		}
		if processor != nil {
			result = append(result, processor)
		}
	}
	return &result, nil
}

// Apply a processor pipeline on an article
func (pp *ProcessorPipeline) Apply(ctx context.Context, article *model.ArticleForm) (bool, error) {
	for _, processor := range *pp {
		applied, err := processor.Apply(ctx, article)
		if err != nil {
			return false, err
		}
		if applied {
			return true, nil
		}
	}
	return false, nil
}
