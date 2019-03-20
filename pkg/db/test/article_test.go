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

func TestGetPaginatedArticlesByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-003")

	// Page request
	req := model.ArticlesPageRequest{
		Limit: 2,
	}

	res, err := testDB.GetPaginatedArticlesByUserID(*user.ID, req)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, res, "response shouldn't be nil")
	assert.True(t, res.TotalCount >= 0, "total count should be a positive integer")
	assert.True(t, !res.HasNext, "we should only have one page")
	assert.True(t, len(res.Entries) >= 0, "entries shouldn't be empty")
}
