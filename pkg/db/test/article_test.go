package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func assertNewArticle(t *testing.T, article *model.Article) *model.Article {
	article, err := testDB.CreateOrUpdateArticle(*article)
	assert.Nil(t, err, "error on create/update article should be nil")
	assert.NotNil(t, article, "article shouldn't be nil")
	assert.NotNil(t, article.ID, "article ID shouldn't be nil")
	return article
}

func TestCreateOrUpdateArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-003")
	// Assert category exists
	category := assertCategoryExists(t, user.ID, "My test category")

	// Create article test case
	builder := model.NewArticleBuilder()
	article := builder.UserID(
		*user.ID,
	).CategoryID(
		*category.ID,
	).Random().Build()

	newArticle := assertNewArticle(t, article)
	assert.Equal(t, article.Title, newArticle.Title, "")
}

func TestDeleteArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-003")
	// Assert category exists
	category := assertCategoryExists(t, user.ID, "My test category")

	// Create article test case
	builder := model.NewArticleBuilder()
	article := builder.UserID(
		*user.ID,
	).CategoryID(
		*category.ID,
	).Random().Build()

	article = assertNewArticle(t, article)

	err := testDB.DeleteArticle(*article)
	assert.Nil(t, err, "error on delete should be nil")

	article, err = testDB.GetArticleByID(*article.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, article == nil, "article should be nil")
}

func TestGetArticlesByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-003")

	articles, err := testDB.GetArticlesByUserID(*user.ID)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, articles, "feed shouldn't be nil")
	assert.True(t, len(articles) >= 0, "articles shouldn't be empty")
}
