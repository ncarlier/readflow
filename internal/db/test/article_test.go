package dbtest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/internal/model"
)

func assertNewArticle(t *testing.T, uid uint, form model.ArticleCreateForm) *model.Article {
	article, err := testDB.CreateArticleForUser(uid, form)
	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.NotNil(t, article.ID)
	assert.Equal(t, form.Title, article.Title)
	assert.Equal(t, "inbox", article.Status, "article status should be inbox")
	assert.Equal(t, uint(0), article.Stars)
	return article
}

func TestCreateAndUpdateArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	category := assertCategoryExists(t, uid, "My test category", "none")

	// Create article test case
	builder := model.NewArticleCreateFormBuilder()
	create := builder.Random().CategoryID(
		*category.ID,
	).Build()

	article := assertNewArticle(t, uid, *create)
	assert.Equal(t, *category.ID, *article.CategoryID)
	updatedAt := *article.UpdatedAt

	// Update article
	status := "read"
	title := article.Title + " (updated)"
	text := "(updated) " + *article.Text
	update := model.ArticleUpdateForm{
		ID:     article.ID,
		Status: &status,
		Title:  &title,
		Text:   &text,
	}
	article, err := testDB.UpdateArticleForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, "read", article.Status, "article status should be read")
	assert.NotEqual(t, updatedAt, *article.UpdatedAt)
	assert.True(t, strings.HasSuffix(article.Title, "(updated)"), "article title should be updated")
	assert.True(t, strings.HasPrefix(*article.Text, "(updated)"), "article text should be updated")

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
	req := model.ArticlesPageRequest{}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, res.TotalCount, uint(0), "total count should be a positive integer")
	assert.False(t, res.HasNext, "we should only have one page")
	assert.GreaterOrEqual(t, len(res.Entries), 0, "entries shouldn't be empty")
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
	status := "inbox"
	req := model.ArticlesPageRequest{
		Status: &status,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Positive(t, res.TotalCount, "total count should be a positive integer")

	nb, err := testDB.MarkAllArticlesAsReadByUser(uid, "inbox", nil)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, nb, "all articles should be marked as read")

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
		Status: &status,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Greater(t, res.TotalCount, uint(0), "total count should be a positive integer")

	nb, err := testDB.DeleteAllReadArticlesByUser(uid)
	assert.Nil(t, err)
	assert.Equal(t, res.TotalCount, uint(nb), "unexpected number of deleted articles")

	res, err = testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(0), res.TotalCount)
}

func TestStarredArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create article test case
	uid := *testUser.ID
	builder := model.NewArticleCreateFormBuilder()
	create := builder.Random().Build()
	article := assertNewArticle(t, uid, *create)

	// Update article
	status := "read"
	stars := 2
	update := model.ArticleUpdateForm{
		ID:     article.ID,
		Status: &status,
		Stars:  &stars,
	}
	article, err := testDB.UpdateArticleForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, "read", article.Status, "article status should be read")
	assert.Equal(t, uint(2), article.Stars, "article status should be starred")

	// Try to delate all read articles
	nb, err := testDB.DeleteAllReadArticlesByUser(uid)
	assert.Nil(t, err)
	assert.Equal(t, uint(0), uint(nb), "unexpected number of deleted articles")

	// Cleanup
	err = testDB.DeleteArticle(article.ID)
	assert.Nil(t, err)
	article, err = testDB.GetArticleByID(article.ID)
	assert.Nil(t, err)
	assert.Nil(t, article)
}
