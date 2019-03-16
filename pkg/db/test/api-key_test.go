package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func assertAPIKeyExists(t *testing.T, userID uint32, alias string) *model.APIKey {
	apiKey, err := testDB.GetAPIKeyByUserIDAndAlias(userID, alias)
	assert.Nil(t, err, "error on getting apiKey by user and alias should be nil")
	if apiKey != nil {
		return apiKey
	}

	builder := model.NewAPIKeyBuilder()
	apiKey = builder.UserID(userID).Alias(alias).Build()

	apiKey, err = testDB.CreateOrUpdateAPIKey(*apiKey)
	assert.Nil(t, err, "error on create/update apiKey should be nil")
	assert.NotNil(t, apiKey, "apiKey shouldn't be nil")
	assert.NotNil(t, apiKey.ID, "apiKey ID shouldn't be nil")
	assert.Equal(t, alias, apiKey.Alias, "")
	assert.True(t, apiKey.Token != "", "apiKey token shouldn't be empty")
	return apiKey
}
func TestCreateOrUpdateAPIKey(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	alias := "My test apiKey"

	// Assert user exists
	user := assertUserExists(t, "test-004")

	assertAPIKeyExists(t, *user.ID, alias)
}

func TestDeleteAPIKey(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	alias := "My apiKey"

	// Assert user exists
	user := assertUserExists(t, "test-004")

	// Assert apiKey exists
	apiKey := assertAPIKeyExists(t, *user.ID, alias)

	err := testDB.DeleteAPIKey(*apiKey)
	assert.Nil(t, err, "error on delete should be nil")

	apiKey, err = testDB.GetAPIKeyByToken(apiKey.Token)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, apiKey == nil, "apiKey should be nil")
}

func TestGetAPIKeysByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-004")

	apiKeys, err := testDB.GetAPIKeysByUserID(*user.ID)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, apiKeys, "apiKeys shouldn't be nil")
	assert.True(t, len(apiKeys) >= 0, "apiKeys shouldn't be empty")
}
