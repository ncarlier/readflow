package ruleengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/model"
)

// RuleProcessor define a rule processor
type RuleProcessor struct {
	category model.Category
	program  *vm.Program
}

// NewRuleProcessor creates a new rule processor
func NewRuleProcessor(category model.Category) (*RuleProcessor, error) {
	if category.Rule == nil || *category.Rule == "" {
		return nil, nil
	}
	env := map[string]interface{}{
		"title":   "",
		"text":    "",
		"url":     "",
		"webhook": "",
		"tags":    []string{},
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}
	p, err := expr.Compile(*category.Rule, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("invalid rule expression: %s", err)
	}
	return &RuleProcessor{
		category: category,
		program:  p,
	}, nil
}

// Apply a rule on an article
func (rp *RuleProcessor) Apply(ctx context.Context, article *model.ArticleCreateForm) (bool, error) {
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
	incomingWebhookAlias := ""
	if alias := ctx.Value(constant.ContextIncomingWebhookAlias); alias != nil {
		incomingWebhookAlias = alias.(string)
	}

	env := map[string]interface{}{
		"title":   article.Title,
		"text":    text,
		"url":     url,
		"webhook": incomingWebhookAlias,
		"tags":    tags,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}

	result, err := expr.Run(rp.program, env)
	if err != nil {
		return false, fmt.Errorf("unable to eval expression on article %s: %s", *article.URL, err)
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
func (pp *ProcessorPipeline) Apply(ctx context.Context, article *model.ArticleCreateForm) (bool, error) {
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
