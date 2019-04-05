package test

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
	ruleengine "github.com/ncarlier/reader/pkg/rule-engine"
)

func newTestRule(rule string, category uint) model.Rule {
	id := uint(1)
	return model.Rule{
		ID:         &id,
		Alias:      "test",
		CategoryID: &category,
		Rule:       rule,
	}
}

func TestBadRuleProcessor(t *testing.T) {
	rule := newTestRule("", 1)
	processor, err := ruleengine.NewRuleProcessor(rule)
	assert.NotNil(t, err, "error should be not nil")
	assert.True(t, processor == nil, "processor should be nil")
}

func TestRuleProcessor(t *testing.T) {
	categoryID := uint(9)
	rule := newTestRule("article.Title == \"test\"", categoryID)
	processor, err := ruleengine.NewRuleProcessor(rule)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, processor != nil, "processor should not be nil")

	builder := model.NewArticleBuilder()
	article := builder.Random().UserID(uint(1)).Title("test").Build()
	applied, err := processor.Apply(article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "processor should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, categoryID, *article.CategoryID, "category should be updated")

	builder = model.NewArticleBuilder()
	article = builder.Random().UserID(uint(1)).Title("foo").Build()
	applied, err = processor.Apply(article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, !applied, "processor should not be applied")
	assert.True(t, article.CategoryID == nil, "category should be nil")
}

func TestProcessorPipeline(t *testing.T) {
	rules := []model.Rule{
		newTestRule("article.Title == \"test\"", uint(1)),
		newTestRule("article.Title == \"foo\"", uint(2)),
		newTestRule("article.Title == \"bar\"", uint(3)),
	}
	pipeline, err := ruleengine.NewProcessorsPipeline(rules)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, pipeline != nil, "pipeline should not be nil")

	builder := model.NewArticleBuilder()
	article := builder.Random().UserID(uint(1)).Title("foo").Build()
	applied, err := pipeline.Apply(article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "pipeline should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, uint(2), *article.CategoryID, "category should be updated")

	builder = model.NewArticleBuilder()
	article = builder.Random().UserID(uint(1)).Title("other").Build()
	applied, err = pipeline.Apply(article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, !applied, "pipeline should not be applied")
	assert.True(t, article.CategoryID == nil, "category should be nil")
}
