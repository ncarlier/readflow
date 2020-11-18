package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertOutboundServiceExists(t *testing.T, uid uint, form model.OutboundServiceCreateForm) *model.OutboundService {
	result, err := testDB.GetOutboundServiceByUserAndAlias(uid, &form.Alias)
	assert.Nil(t, err)
	if result != nil {
		return result
	}

	result, err = testDB.CreateOutboundServiceForUser(uid, form)
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
func TestOutboundServiceCRUD(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create outboundService
	uid := *testUser.ID
	builder := model.NewOutboundServiceCreateFormBuilder()
	create := builder.Alias(
		"My test outbound service",
	).Dummy().Build()
	outboundService := assertOutboundServiceExists(t, uid, *create)

	alias := "My updated outbound service"
	update := model.OutboundServiceUpdateForm{
		ID:    *outboundService.ID,
		Alias: &alias,
	}

	// Update outboundService
	outboundService, err := testDB.UpdateOutboundServiceForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, outboundService)
	assert.NotNil(t, outboundService.ID)
	assert.Equal(t, alias, outboundService.Alias)
	assert.Equal(t, create.Provider, outboundService.Provider)
	assert.Equal(t, create.Config, outboundService.Config)
	assert.Equal(t, create.IsDefault, outboundService.IsDefault)

	// Cleanup
	err = testDB.DeleteOutboundServiceByUser(uid, *outboundService.ID)
	assert.Nil(t, err)
}

func TestDeleteOutboundService(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create outboundService
	uid := *testUser.ID
	builder := model.NewOutboundServiceCreateFormBuilder()
	create := builder.Alias(
		"My test outbound service",
	).Dummy().Build()
	outboundService := assertOutboundServiceExists(t, uid, *create)

	// Retrieve outboundService
	outboundServices, err := testDB.GetOutboundServicesByUser(uid)
	assert.Nil(t, err)
	assert.True(t, len(outboundServices) > 0)

	// Delete outboundService
	err = testDB.DeleteOutboundServiceByUser(uid, *outboundService.ID)
	assert.Nil(t, err)

	// Unable to retrieve deleted outboundService
	outboundService, err = testDB.GetOutboundServiceByUserAndAlias(uid, &create.Alias)
	assert.Nil(t, err)
	assert.Nil(t, outboundService)
}

func TestUpdateDefaultOutboundService(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create first default outboundService
	uid := *testUser.ID
	builder := model.NewOutboundServiceCreateFormBuilder()
	create := builder.Alias(
		"My default outbound service",
	).Dummy().IsDefault(true).Build()
	first := assertOutboundServiceExists(t, uid, *create)

	// Create second default outboundService
	builder = model.NewOutboundServiceCreateFormBuilder()
	create = builder.Alias(
		"My second outbound service",
	).Dummy().IsDefault(true).Build()
	second := assertOutboundServiceExists(t, uid, *create)

	// Refresh first outboundService
	first, err := testDB.GetOutboundServiceByID(*first.ID)
	assert.Nil(t, err)
	assert.False(t, first.IsDefault, "outbound service should not be the default anymore")

	// Test default outboundService query
	defaultOutboundService, err := testDB.GetOutboundServiceByUserAndAlias(uid, nil)
	assert.Nil(t, err)
	assert.NotNil(t, defaultOutboundService)
	assert.NotEqual(t, *first.ID, *defaultOutboundService.ID)
	assert.Equal(t, *second.ID, *defaultOutboundService.ID)

	// Cleanup
	err = testDB.DeleteOutboundServiceByUser(uid, *first.ID)
	assert.Nil(t, err)
	err = testDB.DeleteOutboundServiceByUser(uid, *second.ID)
	assert.Nil(t, err)
}
