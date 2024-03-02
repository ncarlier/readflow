package dbtest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/internal/db"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/logger"
)

var testDB db.DB

const defaultUsername = "test"

const defaultDBConnString = "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"

var testUser *model.User

func assertUserExists(t *testing.T, username string) *model.User {
	user, err := testDB.GetUserByUsername(username)
	assert.Nil(t, err)
	if user != nil {
		return user
	}

	user = &model.User{
		Username: username,
	}
	user, err = testDB.CreateOrUpdateUser(*user)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, user.ID)
	assert.Equal(t, username, user.Username, "")
	return user
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	var err error
	testDB, err = db.NewDB(getEnv("DB", defaultDBConnString))
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	testUser = assertUserExists(t, defaultUsername)
	return func(t *testing.T) {
		t.Log("teardown test case")
		defer testDB.Close()
		// err := testDB.DeleteUser(*testUser)
		// assert.Nil(t, err, "error should be nil")
	}
}

func init() {
	logger.Configure("debug", "text", "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("READFLOW_" + key); ok {
		return value
	}
	return fallback
}
