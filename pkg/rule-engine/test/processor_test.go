package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ncarlier/readflow/pkg/constant"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
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

type testCase struct {
	article  *model.ArticleForm
	rule     string
	category *uint
	err      error
}

var testCases = []testCase{}

var testCategory = uint(99)

func init() {
	builder := model.NewArticleBuilder()
	testCases = append(testCases, testCase{
		article:  builder.Random().BuildForm(""),
		rule:     "",
		category: nil,
		err:      errors.New("invalid rule expression: unexpected token EOF"),
	}, testCase{
		article:  builder.Random().Title("foo").BuildForm(""),
		rule:     "title == \"foo\"",
		category: &testCategory,
		err:      nil,
	}, testCase{
		article:  builder.Random().Text("bar foo bar").BuildForm(""),
		rule:     "text matches \"foo\"",
		category: &testCategory,
		err:      nil,
	}, testCase{
		article:  builder.Random().Text("bar bar bar").BuildForm(""),
		rule:     "text matches \"foo\"",
		category: nil,
		err:      nil,
	}, testCase{
		article:  builder.Random().BuildForm("foo,bar"),
		rule:     "\"foo\" in tags",
		category: &testCategory,
		err:      nil,
	}, testCase{
		article:  builder.Random().BuildForm("foo,bar"),
		rule:     "\"test\" not in tags",
		category: &testCategory,
		err:      nil,
	})
}

func TestRulesTestCases(t *testing.T) {
	ctx := context.TODO()

	for idx, tc := range testCases {
		prefix := fmt.Sprintf("#%d: ", idx)
		rule := newTestRule(tc.rule, testCategory)
		processor, err := ruleengine.NewRuleProcessor(rule)
		if tc.err != nil {
			assert.NotNil(t, err, prefix+"error should be not nil")
			assert.Equal(t, tc.err.Error(), err.Error(), prefix+"error should be equal")
			assert.True(t, processor == nil, prefix+"processor should be nil")
			continue
		}
		assert.Nil(t, err, prefix+"error should be nil")
		assert.True(t, processor != nil, prefix+"processor should not be nil")
		applied, err := processor.Apply(ctx, tc.article)
		assert.Nil(t, err, prefix+"error should be nil")
		if tc.category == nil {
			assert.True(t, !applied, prefix+"processor should not be applied")
		} else {
			assert.True(t, applied, prefix+"processor should be applied")
			assert.True(t, tc.article.CategoryID != nil, prefix+"category should be not nil")
			assert.Equal(t, *tc.category, *tc.article.CategoryID, prefix+"category should be updated")
		}
	}
}

func TestProcessorPipeline(t *testing.T) {
	ctx := context.TODO()
	rules := []model.Rule{
		newTestRule("title == \"test\"", uint(1)),
		newTestRule("text matches \"foo\"", uint(2)),
		newTestRule("\"foo\" in tags", uint(3)),
	}
	pipeline, err := ruleengine.NewProcessorsPipeline(rules)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, pipeline != nil, "pipeline should not be nil")

	builder := model.NewArticleBuilder()
	article := builder.Random().UserID(uint(1)).Text("foo bar foo").BuildForm("")
	applied, err := pipeline.Apply(ctx, article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "pipeline should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, uint(2), *article.CategoryID, "category should be updated")

	builder = model.NewArticleBuilder()
	article = builder.Random().UserID(uint(1)).Text("other").BuildForm("")
	applied, err = pipeline.Apply(ctx, article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, !applied, "pipeline should not be applied")
	assert.True(t, article.CategoryID == nil, "category should be nil")
}
func TestRuleProcessorWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), constant.APIKeyAlias, "test")
	categoryID := uint(9)
	rule := newTestRule("key == \"test\"", categoryID)
	processor, err := ruleengine.NewRuleProcessor(rule)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, processor != nil, "processor should not be nil")

	builder := model.NewArticleBuilder()
	article := builder.Random().UserID(uint(1)).Title("test").BuildForm("")
	applied, err := processor.Apply(ctx, article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "processor should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, categoryID, *article.CategoryID, "category should be updated")
}
