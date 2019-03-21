package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func assertArchiverExists(t *testing.T, archiver model.Archiver) *model.Archiver {
	result, err := testDB.GetArchiverByUserIDAndAlias(*archiver.UserID, archiver.Alias)
	assert.Nil(t, err, "error on getting archiver by user and title should be nil")
	if result != nil {
		return result
	}

	result, err = testDB.CreateOrUpdateArchiver(archiver)
	assert.Nil(t, err, "error on create/update archiver should be nil")
	assert.NotNil(t, result, "archiver shouldn't be nil")
	assert.NotNil(t, result.ID, "archiver ID shouldn't be nil")
	assert.Equal(t, *archiver.UserID, *result.UserID, "")
	assert.Equal(t, archiver.Alias, result.Alias, "")
	assert.Equal(t, archiver.Provider, result.Provider, "")
	return result
}
func TestCreateOrUpdateArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert user exists
	user := assertUserExists(t, "test-005")

	archiver := model.Archiver{
		Alias:    "My archiver",
		UserID:   user.ID,
		Provider: "test",
		Config:   "{\"foo\": \"bar\"}",
	}

	// Create archiver
	newArchiver := assertArchiverExists(t, archiver)

	newArchiver.Alias = "My updated archiver"

	// Update archiver
	updatedArchiver, err := testDB.CreateOrUpdateArchiver(*newArchiver)
	assert.Nil(t, err, "error on update archiver should be nil")
	assert.NotNil(t, updatedArchiver, "archiver shouldn't be nil")
	assert.NotNil(t, updatedArchiver.ID, "archiver ID shouldn't be nil")
	assert.Equal(t, newArchiver.Alias, updatedArchiver.Alias, "")
}

func TestDeleteArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	user := assertUserExists(t, "test-002")

	// Assert archiver exists
	archiver := &model.Archiver{
		Alias:    "My archiver",
		UserID:   user.ID,
		Provider: "test",
		Config:   "{\"foo\": \"bar\"}",
	}
	archiver = assertArchiverExists(t, *archiver)

	archivers, err := testDB.GetArchiversByUserID(*user.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(archivers) > 0, "archivers should not be empty")

	err = testDB.DeleteArchiver(*archiver)
	assert.Nil(t, err, "error on delete should be nil")

	archiver, err = testDB.GetArchiverByUserIDAndAlias(*user.ID, "My archiver")
	assert.Nil(t, err, "error should be nil")
	assert.True(t, archiver == nil, "archiver should be nil")
}
