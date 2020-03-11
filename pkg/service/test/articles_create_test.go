package dbtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestCreateArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	url := "https://github.com/ncarlier/readflow"
	req := model.ArticleForm{
		URL: &url,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err, "")
	assert.Equal(t, "ncarlier/readflow", art.Title, "")
	assert.NotNil(t, art.URL, "")
	assert.Equal(t, url, *art.URL, "")
	assert.NotNil(t, art.Image, "")
	assert.True(t, strings.HasPrefix(*art.Image, "https://repository-images.githubusercontent.com"), "")
	assert.True(t, art.CategoryID == nil, "")

	// Create same article again
	_, err = service.Lookup().CreateArticle(testContext, req, opts)
	assert.Equal(t, "already exists", err.Error(), "")
}

func TestCreateArticleInCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	categoryTitle := "test category"
	cat, err := service.Lookup().CreateOrUpdateCategory(testContext, nil, categoryTitle)
	assert.Nil(t, err, "")
	assert.Equal(t, categoryTitle, cat.Title, "")
	assert.Equal(t, *testUser.ID, *cat.UserID, "")

	// Create article
	req := model.ArticleForm{
		Title:      "TestCreateArticleInCategory",
		CategoryID: cat.ID,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err, "")
	assert.Equal(t, req.Title, art.Title, "")
	assert.True(t, art.CategoryID != nil, "")
	assert.Equal(t, *cat.ID, *art.CategoryID, "")
}

func TestCreateArticleWithRuleEngine(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	categoryTitle := "test category"
	cat, err := service.Lookup().CreateOrUpdateCategory(testContext, nil, categoryTitle)
	assert.Nil(t, err, "")
	assert.Equal(t, categoryTitle, cat.Title, "")
	assert.Equal(t, *testUser.ID, *cat.UserID, "")

	// Create rule
	form := model.RuleForm{
		Alias:      "test-rule",
		CategoryID: *cat.ID,
		Priority:   1,
		Rule:       "title matches \"^Test\"",
	}
	rule, err := service.Lookup().CreateOrUpdateRule(testContext, form)
	assert.Nil(t, err, "")
	assert.Equal(t, form.Alias, rule.Alias, "")

	// Create article
	req := model.ArticleForm{
		Title: "TestCreateArticleWithRuleEngine",
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err, "")
	assert.Equal(t, req.Title, art.Title, "")
	assert.True(t, art.CategoryID != nil, "")
	assert.Equal(t, *cat.ID, *art.CategoryID, "")
}

func TestCreateArticlesExceedingQuota(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create 6 articles (quota is 5)
	for i := 1; i <= 6; i++ {
		req := model.ArticleForm{
			Title: fmt.Sprintf("TestCreateArticlesExceedingQuota-%d", i),
		}
		opts := service.ArticleCreationOptions{}
		_, err := service.Lookup().CreateArticle(testContext, req, opts)
		if i <= 5 {
			assert.Nil(t, err, "")
		} else {
			assert.Equal(t, service.ErrUserQuotaReached.Error(), err.Error(), "")
		}
	}
}
