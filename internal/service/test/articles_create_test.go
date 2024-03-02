package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
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
	require.Nil(t, err)
	require.Contains(t, art.Title, "GitHub - ncarlier/readflow")
	require.NotNil(t, art.URL)
	require.Equal(t, url, *art.URL)
	require.NotNil(t, art.Image)
	require.True(t, strings.HasPrefix(*art.Image, "https://repository-images.githubusercontent.com"), "unexpected image URL")
	require.Nil(t, art.CategoryID)

	// Create same article again
	_, err = service.Lookup().CreateArticle(testContext, req, opts)
	require.Equal(t, "already exists", err.Error())
}

func TestCreateArticleInCategory(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	cat := requireNewCategory(t)

	// Create article
	req := model.ArticleCreateForm{
		Title:      "TestCreateArticleInCategory",
		CategoryID: cat.ID,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	require.Nil(t, err)
	require.Equal(t, req.Title, art.Title)
	require.NotNil(t, art.CategoryID)
	require.Equal(t, *cat.ID, *art.CategoryID)
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
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
			require.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
