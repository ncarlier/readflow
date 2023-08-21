package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestCreateRemoteArticle(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	url := "https://github.com/ncarlier/readflow"
	req := model.ArticleCreateForm{
		URL: &url,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err)
	assert.Contains(t, art.Title, "GitHub - ncarlier/readflow")
	assert.NotNil(t, art.URL)
	assert.Equal(t, url, *art.URL)
	assert.NotNil(t, art.Image)
	assert.True(t, strings.HasPrefix(*art.Image, "https://repository-images.githubusercontent.com"), "unexpected image URL")
	assert.Nil(t, art.CategoryID)

	// Create same article again
	_, err = service.Lookup().CreateArticle(testContext, req, opts)
	assert.Equal(t, "already exists", err.Error())
}

func TestCreateArticleInCategory(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	cat := assertNewCategory(t)

	// Create article
	req := model.ArticleCreateForm{
		Title:      "TestCreateArticleInCategory",
		CategoryID: cat.ID,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err)
	assert.Equal(t, req.Title, art.Title)
	assert.NotNil(t, art.CategoryID)
	assert.Equal(t, *cat.ID, *art.CategoryID)
}

func TestCreateArticlesExceedingQuota(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create 6 articles (quota is 5)
	for i := 1; i <= 6; i++ {
		req := model.ArticleCreateForm{
			Title: fmt.Sprintf("TestCreateArticlesExceedingQuota-%d", i),
		}
		opts := service.ArticleCreationOptions{}
		_, err := service.Lookup().CreateArticle(testContext, req, opts)
		if i <= 5 {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
