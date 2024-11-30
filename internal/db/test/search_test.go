package dbtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
)

func TestSearchArticlesByFullTextQuery(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	category := assertCategoryExists(t, uid, "My test category", "none")

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

	status := "inbox"
	query := "science +computer"
	// Page request
	req := model.ArticlesPageRequest{
		Status: &status,
		Query:  &query,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, uint(1), res.TotalCount)
	require.False(t, res.HasNext, "we should only have one page")
	require.Equal(t, 1, len(res.Entries))
	require.Equal(t, "About computer science", res.Entries[0].Title)
}

func TestSearchArticlesByHashtags(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID

	// Create articles test case
	builder := model.NewArticleCreateFormBuilder()
	create := builder.Random().Title("About #cool subject").Build()
	assertNewArticle(t, uid, *create)
	create = builder.Random().Title("About cool #subject").Build()
	assertNewArticle(t, uid, *create)

	status := "inbox"
	query := "#cool"
	// Page request
	req := model.ArticlesPageRequest{
		Status: &status,
		Query:  &query,
	}

	res, err := testDB.GetPaginatedArticlesByUser(uid, req)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, uint(1), res.TotalCount)
	require.Equal(t, 1, len(res.Entries))
	require.Equal(t, "About #cool subject", res.Entries[0].Title)

}
