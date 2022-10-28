package test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestUpdateArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := assertNewArticle(t, "article update test")

	// Create category
	cat := assertNewCategory(t)

	title := "readflow Github project page"
	update := model.ArticleUpdateForm{
		ID:         art.ID,
		Title:      &title,
		CategoryID: cat.ID,
	}
	// Update article
	art, err := service.Lookup().UpdateArticle(testContext, update)
	assert.Nil(t, err)
	assert.Equal(t, title, art.Title)
	assert.NotNil(t, art.CategoryID)
	assert.Equal(t, *cat.ID, *art.CategoryID)
}

func TestUpdateArticleWithErrors(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := assertNewArticle(t, "article update test with errors")

	title := gofakeit.Sentence(99)
	stars := uint(10)
	update := model.ArticleUpdateForm{
		ID:    art.ID,
		Title: &title,
		Stars: &stars,
	}
	// Update article
	_, err := service.Lookup().UpdateArticle(testContext, update)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "title")
	assert.Contains(t, err.Error(), "stars")
}

func TestUpdateArticleWithBadCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	art := assertNewArticle(t, "article update test with bad category")

	catID := uint(999)
	update := model.ArticleUpdateForm{
		ID:         art.ID,
		CategoryID: &catID,
	}
	// Update article
	art, err := service.Lookup().UpdateArticle(testContext, update)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "category not found")
}
