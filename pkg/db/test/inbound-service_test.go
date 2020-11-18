package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertInboundServiceExists(t *testing.T, uid uint, alias string) *model.InboundService {
	inboundService, err := testDB.GetInboundServiceByUserAndAlias(uid, alias)
	assert.Nil(t, err)
	if inboundService != nil {
		return inboundService
	}

	builder := model.NewInboundServiceCreateFormBuilder()
	form := builder.Alias(alias).Dummy().Build()

	inboundService, err = testDB.CreateInboundServiceForUser(uid, *form)
	assert.Nil(t, err)
	assert.NotNil(t, inboundService)
	assert.NotNil(t, inboundService.ID)
	assert.Equal(t, alias, inboundService.Alias)
	assert.NotEqual(t, "", inboundService.Token)
	assert.Equal(t, "dummy", inboundService.Provider)
	return inboundService
}
func TestCreateOrUpdateInboundService(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My test inbound service"

	assertInboundServiceExists(t, uid, alias)
}

func TestDeleteInboundService(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My inbound service"

	// Assert inboundService exists
	inboundService := assertInboundServiceExists(t, uid, alias)

	err := testDB.DeleteInboundServiceByUser(uid, *inboundService.ID)
	assert.Nil(t, err)

	inboundService, err = testDB.GetInboundServiceByToken(inboundService.Token)
	assert.Nil(t, err)
	assert.Nil(t, inboundService)
}

func TestGetInboundServicesByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID

	inboundServices, err := testDB.GetInboundServicesByUser(uid)
	assert.Nil(t, err)
	assert.NotNil(t, inboundServices)
	assert.True(t, len(inboundServices) >= 0)
}
