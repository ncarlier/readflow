package db_test

import (
	"testing"

	"github.com/ncarlier/reader/pkg/logger"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/db"
)

var testDB db.DB

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	var err error
	testDB, err = db.Configure("postgres://postgres:testpwd@localhost/reader_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to setup Database: %v", err)
	}
	return func(t *testing.T) {
		t.Log("teardown test case")
		testDB.Close()
	}
}

func TestGetArticles(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	articles, err := testDB.GetArticles()
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, articles, "feed shouldn't be nil")
	assert.True(t, len(articles) >= 0, "articles shouldn't be empty")
}

func init() {
	logger.Configure("debug", true, nil)
}
