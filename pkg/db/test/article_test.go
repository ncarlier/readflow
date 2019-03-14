package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
)

func TestGetArticles(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	articles, err := testDB.GetArticles()
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, articles, "feed shouldn't be nil")
	assert.True(t, len(articles) >= 0, "articles shouldn't be empty")
}
