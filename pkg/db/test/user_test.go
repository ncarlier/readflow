package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func TestCreateOrUpdateUser(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	username := "test-001"

	user := &model.User{
		Username: username,
	}
	user, err := testDB.CreateOrUpdateUser(*user)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, user, "user shouldn't be nil")
	assert.NotNil(t, user.ID, "user ID shouldn't be nil")
	assert.Equal(t, username, user.Username, "")
	assert.True(t, *user.ID > 0, "user ID should be a valid integer")
	assert.True(t, !user.Enabled, "user should be disabled")
}

func TestDeleteUser(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	username := "test-001"

	user, err := testDB.GetUserByUsername(username)
	assert.Nil(t, err, "error should be nil")
	if user == nil {
		user = &model.User{
			Username: username,
		}
		_, err = testDB.CreateOrUpdateUser(*user)
		assert.Nil(t, err, "error should be nil")
	}

	user = &model.User{
		Username: username,
	}
	err = testDB.DeleteUser(*user)
	assert.Nil(t, err, "error should be nil")

	user, err = testDB.GetUserByUsername(username)
	assert.Nil(t, err, "error should be nil")
	// assert.Nil(t, user, "user should be nil") // NOT WORKING
	assert.True(t, user == nil, "user should be nil")
}
