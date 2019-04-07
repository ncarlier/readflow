package dbtest

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
)

func assertArchiverExists(t *testing.T, archiver model.Archiver) *model.Archiver {
	result, err := testDB.GetArchiverByUserIDAndAlias(*archiver.UserID, &archiver.Alias)
	assert.Nil(t, err, "error on getting archiver by user and alias should be nil")
	if result != nil {
		id := *result.ID
		archiver.ID = &id
	}

	result, err = testDB.CreateOrUpdateArchiver(archiver)
	assert.Nil(t, err, "error on create/update archiver should be nil")
	assert.NotNil(t, result, "archiver shouldn't be nil")
	assert.NotNil(t, result.ID, "archiver ID shouldn't be nil")
	assert.Equal(t, *archiver.UserID, *result.UserID, "")
	assert.Equal(t, archiver.Alias, result.Alias, "")
	assert.Equal(t, archiver.Provider, result.Provider, "")
	assert.Equal(t, archiver.IsDefault, result.IsDefault, "")
	return result
}
func TestCreateOrUpdateArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	archiver := model.Archiver{
		Alias:    "My archiver",
		UserID:   testUser.ID,
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

	// Cleanup
	err = testDB.DeleteArchiver(*updatedArchiver)
	assert.Nil(t, err, "error on cleanup should be nil")
}

func TestDeleteArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	alias := "My archiver"

	// Assert archiver exists
	archiver := &model.Archiver{
		Alias:    alias,
		UserID:   testUser.ID,
		Provider: "test",
		Config:   "{\"foo\": \"bar\"}",
	}
	archiver = assertArchiverExists(t, *archiver)

	archivers, err := testDB.GetArchiversByUserID(*testUser.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(archivers) > 0, "archivers should not be empty")

	err = testDB.DeleteArchiver(*archiver)
	assert.Nil(t, err, "error on delete should be nil")

	archiver, err = testDB.GetArchiverByUserIDAndAlias(*testUser.ID, &alias)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, archiver == nil, "archiver should be nil")
}

func TestUpdateDefaultArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert archiver exists
	archiver := model.Archiver{
		Alias:     "My default archiver",
		UserID:    testUser.ID,
		Provider:  "test",
		Config:    "{\"foo\": \"bar\"}",
		IsDefault: true,
	}
	firstArchiver := assertArchiverExists(t, archiver)
	archiver.Alias = "My new default archiver"
	assertArchiverExists(t, archiver)

	// Refresh first archiver
	firstArchiver, err := testDB.GetArchiverByID(*firstArchiver.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, !firstArchiver.IsDefault, "archiver should not be the default anymore")

	// Test default archiver query
	defaultArchiver, err := testDB.GetArchiverByUserIDAndAlias(*testUser.ID, nil)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, defaultArchiver != nil, "default archiver should noty be nil")
	assert.NotEqual(t, defaultArchiver.ID, firstArchiver.ID, "")
}
