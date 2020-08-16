package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func TestSearchArticlesByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	category := assertCategoryExists(t, uid, "My test category")

	// Create articles test case
	builder := model.NewArticleCreateFormBuilder()
	create := builder.Random().Title("About computer science").CategoryID(
		*category.ID,
	).Build()
	assertNewArticle(t, uid, *create)
	create = builder.Random().Title("About science").CategoryID(
		*category.ID,
	).Build()
	assertNewArticle(t, uid, *create)

	status := "unread"
	query := "science and computer"
	// Page request
	req := model.ArticlesPageRequest{
		Status: &status,
		Query:  &query,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(1), res.TotalCount)
	assert.False(t, res.HasNext, "we should only have one page")
	assert.Equal(t, 1, len(res.Entries))
}
