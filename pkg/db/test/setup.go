package dbtest

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/model"
)

var testDB db.DB

const defaultUsername = "test"

const defaultDBConnString = "postgres://postgres:testpwd@localhost/reader_test?sslmode=disable"

var testUser *model.User

func assertUserExists(t *testing.T, username string) *model.User {
	user, err := testDB.GetUserByUsername(username)
	assert.Nil(t, err, "error getting user by username should be nil")
	if user != nil {
		return user
	}

	user = &model.User{
		Username: username,
	}
	user, err = testDB.CreateOrUpdateUser(*user)
	assert.Nil(t, err, "error on create/update user should be nil")
	assert.NotNil(t, user, "user shouldn't be nil")
	assert.NotNil(t, user.ID, "user ID shouldn't be nil")
	assert.Equal(t, username, user.Username, "")
	return user
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	var err error
	testDB, err = db.Configure(defaultDBConnString)
	if err != nil {
		t.Fatalf("Unable to setup Database: %v", err)
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
	logger.Configure("debug", true, "")
}
