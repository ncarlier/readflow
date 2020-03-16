package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ncarlier/readflow/pkg/constant"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
)

func newTestCategory(id uint, rule string) model.Category {
	builder := model.NewCategoryBuilder()
	category := builder.Random().Rule(rule).Build()
	category.ID = &id
	return *category
}

type testCase struct {
	article       *model.ArticleForm
	category      model.Category
	expectedValue *uint
	expectedError error
}

var testCases = []testCase{}

func init() {
	expectedCategoryID := uint(42)
	builder := model.NewArticleBuilder()
	testCases = append(testCases, testCase{
		article:       builder.Random().BuildForm(""),
		category:      newTestCategory(42, "."),
		expectedValue: nil,
		expectedError: errors.New("invalid rule expression:"),
	}, testCase{
		article:       builder.Random().Title("foo").BuildForm(""),
		category:      newTestCategory(42, "title == \"foo\""),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Text("bar foo bar").BuildForm(""),
		category:      newTestCategory(42, "text matches \"foo\""),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Text("bar bar bar").BuildForm(""),
		category:      newTestCategory(42, "text matches \"foo\""),
		expectedValue: nil,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().BuildForm("foo,bar"),
		category:      newTestCategory(42, "\"foo\" in tags"),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().BuildForm("foo,bar"),
		category:      newTestCategory(42, "\"test\" not in tags"),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	})
}

func TestRulesTestCases(t *testing.T) {
	ctx := context.TODO()

	for idx, tc := range testCases {
		prefix := fmt.Sprintf("#%d: ", idx)
		processor, err := ruleengine.NewRuleProcessor(tc.category)
		if tc.expectedError != nil {
			assert.NotNil(t, err, prefix+"error should be not nil")
			assert.True(t, strings.HasPrefix(err.Error(), tc.expectedError.Error()), prefix+"error should be equal")
			assert.True(t, processor == nil, prefix+"processor should be nil")
			continue
		}
		assert.Nil(t, err, prefix+"error should be nil")
		assert.True(t, processor != nil, prefix+"processor should not be nil")
		applied, err := processor.Apply(ctx, tc.article)
		assert.Nil(t, err, prefix+"error should be nil")
		if tc.expectedValue == nil {
			assert.True(t, !applied, prefix+"processor should not be applied")
		} else {
			assert.True(t, applied, prefix+"processor should be applied")
			assert.True(t, tc.article.CategoryID != nil, prefix+"category should be not nil")
			assert.Equal(t, *tc.expectedValue, *tc.article.CategoryID, prefix+"category should be updated")
		}
	}
}

func TestProcessorPipeline(t *testing.T) {
	ctx := context.TODO()
	categories := []model.Category{
		newTestCategory(1, "title == \"test\""),
		newTestCategory(2, "text matches \"foo\""),
		newTestCategory(3, "\"foo\" in tags"),
	}
	pipeline, err := ruleengine.NewProcessorsPipeline(categories)
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
	category := newTestCategory(9, "key == \"test\"")
	processor, err := ruleengine.NewRuleProcessor(category)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, processor != nil, "processor should not be nil")

	builder := model.NewArticleBuilder()
	article := builder.Random().UserID(uint(1)).Title("test").BuildForm("")
	applied, err := processor.Apply(ctx, article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "processor should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, uint(9), *article.CategoryID, "category should be updated")
}
