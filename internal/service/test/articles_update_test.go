package test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
)

func TestUpdateArticle(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := requireNewArticle(t, "article update test")

	// Create category
	cat := requireNewCategory(t)

	title := "readflow Github project page"
	update := model.ArticleUpdateForm{
		ID:         art.ID,
		Title:      &title,
		CategoryID: cat.ID,
	}
	// Update article
	art, err := service.Lookup().UpdateArticle(testContext, update)
	require.Nil(t, err)
	require.Equal(t, title, art.Title)
	require.NotNil(t, art.CategoryID)
	require.Equal(t, *cat.ID, *art.CategoryID)
}

func TestUpdateArticleWithErrors(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := requireNewArticle(t, "article update test with errors")

	title := gofakeit.Sentence(99)
	stars := 10
	update := model.ArticleUpdateForm{
		ID:    art.ID,
		Title: &title,
		Stars: &stars,
	}
	// Update article
	_, err := service.Lookup().UpdateArticle(testContext, update)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "title")
	require.Contains(t, err.Error(), "stars")
}

func TestUpdateArticleWithBadCategory(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := requireNewArticle(t, "article update test with bad category")

	catID := uint(999)
	update := model.ArticleUpdateForm{
		ID:         art.ID,
		CategoryID: &catID,
	}
	// Update article
	_, err := service.Lookup().UpdateArticle(testContext, update)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "category not found")
}
