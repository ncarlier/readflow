package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertArchiverExists(t *testing.T, uid uint, form model.ArchiverCreateForm) *model.Archiver {
	result, err := testDB.GetArchiverByUserAndAlias(uid, &form.Alias)
	assert.Nil(t, err)
	if result != nil {
		return result
	}

	result, err = testDB.CreateArchiverForUser(uid, form)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ID)
	assert.Equal(t, uid, *result.UserID)
	assert.Equal(t, form.Alias, result.Alias)
	assert.Equal(t, form.Provider, result.Provider)
	assert.Equal(t, form.Config, result.Config)
	assert.Equal(t, form.IsDefault, result.IsDefault)
	return result
}
func TestArchiverCRUD(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create archiver
	uid := *testUser.ID
	builder := model.NewArchiverCreateFormBuilder()
	create := builder.Alias(
		"My test archiver",
	).Dummy().Build()
	archiver := assertArchiverExists(t, uid, *create)

	alias := "My updated archiver"
	update := model.ArchiverUpdateForm{
		ID:    *archiver.ID,
		Alias: &alias,
	}

	// Update archiver
	archiver, err := testDB.UpdateArchiverForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, archiver)
	assert.NotNil(t, archiver.ID)
	assert.Equal(t, alias, archiver.Alias)
	assert.Equal(t, create.Provider, archiver.Provider)
	assert.Equal(t, create.Config, archiver.Config)
	assert.Equal(t, create.IsDefault, archiver.IsDefault)

	// Cleanup
	err = testDB.DeleteArchiverByUser(uid, *archiver.ID)
	assert.Nil(t, err)
}

func TestDeleteArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create archiver
	uid := *testUser.ID
	builder := model.NewArchiverCreateFormBuilder()
	create := builder.Alias(
		"My test archiver",
	).Dummy().Build()
	archiver := assertArchiverExists(t, uid, *create)

	// Retrieve archiver
	archivers, err := testDB.GetArchiversByUser(uid)
	assert.Nil(t, err)
	assert.True(t, len(archivers) > 0)

	// Delete archiver
	err = testDB.DeleteArchiverByUser(uid, *archiver.ID)
	assert.Nil(t, err)

	// Unable to retrieve deleted archiver
	archiver, err = testDB.GetArchiverByUserAndAlias(uid, &create.Alias)
	assert.Nil(t, err)
	assert.Nil(t, archiver)
}

func TestUpdateDefaultArchiver(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create first default archiver
	uid := *testUser.ID
	builder := model.NewArchiverCreateFormBuilder()
	create := builder.Alias(
		"My default archiver",
	).Dummy().IsDefault(true).Build()
	first := assertArchiverExists(t, uid, *create)

	// Create second default archiver
	builder = model.NewArchiverCreateFormBuilder()
	create = builder.Alias(
		"My second archiver",
	).Dummy().IsDefault(true).Build()
	second := assertArchiverExists(t, uid, *create)

	// Refresh first archiver
	first, err := testDB.GetArchiverByID(*first.ID)
	assert.Nil(t, err)
	assert.False(t, first.IsDefault, "archiver should not be the default anymore")

	// Test default archiver query
	defaultArchiver, err := testDB.GetArchiverByUserAndAlias(uid, nil)
	assert.Nil(t, err)
	assert.NotNil(t, defaultArchiver)
	assert.NotEqual(t, *first.ID, *defaultArchiver.ID)
	assert.Equal(t, *second.ID, *defaultArchiver.ID)

	// Cleanup
	err = testDB.DeleteArchiverByUser(uid, *first.ID)
	assert.Nil(t, err)
	err = testDB.DeleteArchiverByUser(uid, *second.ID)
	assert.Nil(t, err)
}
