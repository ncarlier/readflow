package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateUser(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	user := assertUserExists(t, "test-001")
	assert.Positive(t, *user.ID, "user ID should be a valid integer")
	assert.True(t, !user.Enabled, "user should be disabled")

	user.Enabled = true
	user.Plan = "test"
	user, err := testDB.CreateOrUpdateUser(*user)
	assert.Nil(t, err)
	assert.True(t, user != nil)
	assert.True(t, user.Enabled, "user should be enabled")
	assert.Equal(t, "test", user.Plan, "unexpected user plan")
}

func TestDeleteUser(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	username := "test-001"
	user := assertUserExists(t, username)

	err := testDB.DeleteUser(*user)
	assert.Nil(t, err)

	user, err = testDB.GetUserByUsername(username)
	assert.Nil(t, err)
	assert.Nil(t, user)
}
