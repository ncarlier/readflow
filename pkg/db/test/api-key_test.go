package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertAPIKeyExists(t *testing.T, uid uint, alias string) *model.APIKey {
	apiKey, err := testDB.GetAPIKeyByUserAndAlias(uid, alias)
	assert.Nil(t, err)
	if apiKey != nil {
		return apiKey
	}

	builder := model.NewAPIKeyCreateFormBuilder()
	form := builder.Alias(alias).Build()

	apiKey, err = testDB.CreateAPIKeyForUser(uid, *form)
	assert.Nil(t, err)
	assert.NotNil(t, apiKey)
	assert.NotNil(t, apiKey.ID)
	assert.Equal(t, alias, apiKey.Alias)
	assert.NotEqual(t, "", apiKey.Token)
	return apiKey
}
func TestCreateOrUpdateAPIKey(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My test apiKey"

	assertAPIKeyExists(t, uid, alias)
}

func TestDeleteAPIKey(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My apiKey"

	// Assert apiKey exists
	apiKey := assertAPIKeyExists(t, uid, alias)

	err := testDB.DeleteAPIKeyByUser(uid, *apiKey.ID)
	assert.Nil(t, err)

	apiKey, err = testDB.GetAPIKeyByToken(apiKey.Token)
	assert.Nil(t, err)
	assert.Nil(t, apiKey)
}

func TestGetAPIKeysByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID

	apiKeys, err := testDB.GetAPIKeysByUser(uid)
	assert.Nil(t, err)
	assert.NotNil(t, apiKeys)
	assert.True(t, len(apiKeys) >= 0)
}
