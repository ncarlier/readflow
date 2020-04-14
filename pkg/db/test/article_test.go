package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertNewArticle(t *testing.T, uid uint, form model.ArticleCreateForm) *model.Article {
	article, err := testDB.CreateArticleForUser(uid, form)
	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.NotNil(t, article.ID)
	return article
}

func TestCreateAndUpdateArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	category := assertCategoryExists(t, uid, "My test category")

	// Create article test case
	builder := model.NewArticleCreateFormBuilder()
	create := builder.CategoryID(
		*category.ID,
	).Random().Build()

	article := assertNewArticle(t, uid, *create)
	assert.Equal(t, create.Title, article.Title)
	assert.Equal(t, "unread", article.Status, "article status should be unread")
	updatedAt := *article.UpdatedAt

	// Update article
	status := "read"
	update := model.ArticleUpdateForm{
		ID:     article.ID,
		Status: &status,
	}
	article, err := testDB.UpdateArticleForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, "read", article.Status, "article status should be read")
	assert.NotEqual(t, updatedAt, *article.UpdatedAt)

	// Cleanup
	err = testDB.DeleteArticle(article.ID)
	assert.Nil(t, err)
	article, err = testDB.GetArticleByID(article.ID)
	assert.Nil(t, err)
	assert.Nil(t, article)
}

func TestGetPaginatedArticlesByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	// Page request
	req := model.ArticlesPageRequest{
		Limit: 20,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.TotalCount >= 0, "total count should be a positive integer")
	assert.False(t, res.HasNext, "we should only have one page")
	assert.True(t, len(res.Entries) >= 0, "entries shouldn't be empty")
}

func TestMarkAllArticlesAsRead(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create article test case
	uid := *testUser.ID
	builder := model.NewArticleCreateFormBuilder()
	form := builder.Random().Build()
	assertNewArticle(t, uid, *form)

	// Page request
	status := "unread"
	req := model.ArticlesPageRequest{
		Limit:  20,
		Status: &status,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.TotalCount >= 0, "total count should be a positive integer")

	nb, err := testDB.MarkAllArticlesAsReadByUser(uid, nil)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, nb, "all articles sould be marked as read")

	res, err = testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(0), res.TotalCount)
}

func TestDeleteAllReadArticles(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	status := "read"
	req := model.ArticlesPageRequest{
		Limit:  20,
		Status: &status,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.TotalCount >= 0, "total count should be a positive integer")

	nb, err := testDB.DeleteAllReadArticlesByUser(uid)
	assert.Nil(t, err)
	assert.Equal(t, res.TotalCount, uint(nb), "unexpected number of deleted articles")

	res, err = testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(0), res.TotalCount)
}
