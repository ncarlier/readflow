package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
)

func newTestCategory(id uint, rule string) model.Category {
	return model.Category{
		ID:    &id,
		Title: "Dummy category",
		Rule:  &rule,
	}
}

type testCase struct {
	article       *model.ArticleCreateForm
	category      model.Category
	expectedValue *uint
	expectedError error
}

var testCases = []testCase{}

func init() {
	expectedCategoryID := uint(42)
	builder := model.NewArticleCreateFormBuilder()
	testCases = append(testCases, testCase{
		article:       builder.Random().Build(),
		category:      newTestCategory(42, "."),
		expectedValue: nil,
		expectedError: errors.New("invalid rule expression:"),
	}, testCase{
		article:       builder.Random().Title("foo").Build(),
		category:      newTestCategory(42, "title == \"foo\""),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Text("bar foo bar").Build(),
		category:      newTestCategory(42, "text matches \"foo\""),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Text("bar bar bar").Build(),
		category:      newTestCategory(42, "text matches \"foo\""),
		expectedValue: nil,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Tags("foo,bar").Build(),
		category:      newTestCategory(42, "\"foo\" in tags"),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	}, testCase{
		article:       builder.Random().Tags("foo,bar").Build(),
		category:      newTestCategory(42, "\"test\" not in tags"),
		expectedValue: &expectedCategoryID,
		expectedError: nil,
	})
}

func TestRulesTestCases(t *testing.T) {
	ctx := context.TODO()

	for idx, tc := range testCases {
		testCaseName := fmt.Sprintf("Test case #%d", idx)
		processor, err := ruleengine.NewRuleProcessor(tc.category)
		if tc.expectedError != nil {
			assert.NotNil(t, err, testCaseName)
			assert.True(t, strings.HasPrefix(err.Error(), tc.expectedError.Error()), testCaseName)
			assert.Nil(t, processor, testCaseName)
			continue
		}
		assert.Nil(t, err, testCaseName)
		assert.NotNil(t, processor, testCaseName)
		applied, err := processor.Apply(ctx, tc.article)
		assert.Nil(t, err, testCaseName)
		if tc.expectedValue == nil {
			assert.False(t, applied, testCaseName)
		} else {
			assert.True(t, applied, testCaseName)
			assert.NotNil(t, tc.article.CategoryID, testCaseName)
			assert.Equal(t, *tc.expectedValue, *tc.article.CategoryID, testCaseName)
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
	assert.Nil(t, err)
	assert.NotNil(t, pipeline)

	builder := model.NewArticleCreateFormBuilder()
	article := builder.Random().Text("foo bar foo").Build()
	applied, err := pipeline.Apply(ctx, article)
	assert.Nil(t, err)
	assert.True(t, applied)
	assert.NotNil(t, article.CategoryID)
	assert.Equal(t, uint(2), *article.CategoryID)

	builder = model.NewArticleCreateFormBuilder()
	article = builder.Random().Text("other").Build()
	applied, err = pipeline.Apply(ctx, article)
	assert.Nil(t, err)
	assert.False(t, applied)
	assert.Nil(t, article.CategoryID)
}

func TestRuleProcessorWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), constant.APIKeyAlias, "test")
	category := newTestCategory(9, "key == \"test\"")
	processor, err := ruleengine.NewRuleProcessor(category)
	assert.Nil(t, err)
	assert.NotNil(t, processor)

	builder := model.NewArticleCreateFormBuilder()
	article := builder.Random().Title("test").Build()
	applied, err := processor.Apply(ctx, article)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, applied, "processor should be applied")
	assert.True(t, article.CategoryID != nil, "category should be not nil")
	assert.Equal(t, uint(9), *article.CategoryID, "category should be updated")
}
