package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/db"
	"github.com/ncarlier/reader/pkg/logger"
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

func init() {
	logger.Configure("debug", true, nil)
}
