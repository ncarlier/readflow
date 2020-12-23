package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertOutgoingWebhookExists(t *testing.T, uid uint, form model.OutgoingWebhookCreateForm) *model.OutgoingWebhook {
	result, err := testDB.GetOutgoingWebhookByUserAndAlias(uid, &form.Alias)
	assert.Nil(t, err)
	if result != nil {
		return result
	}

	result, err = testDB.CreateOutgoingWebhookForUser(uid, form)
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
func TestOutgoingWebhookCRUD(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create webhook
	uid := *testUser.ID
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	create := builder.Alias(
		"My test outgoing webhook",
	).Dummy().Build()
	webhook := assertOutgoingWebhookExists(t, uid, *create)

	alias := "My updated outgoing webhook"
	update := model.OutgoingWebhookUpdateForm{
		ID:    *webhook.ID,
		Alias: &alias,
	}

	// Update webhook
	webhook, err := testDB.UpdateOutgoingWebhookForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, webhook)
	assert.NotNil(t, webhook.ID)
	assert.Equal(t, alias, webhook.Alias)
	assert.Equal(t, create.Provider, webhook.Provider)
	assert.Equal(t, create.Config, webhook.Config)
	assert.Equal(t, create.IsDefault, webhook.IsDefault)

	// Cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	assert.Nil(t, err)
}

func TestDeleteOutgoingWebhook(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create webhook
	uid := *testUser.ID
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	create := builder.Alias(
		"My test outgoing webhook",
	).Dummy().Build()
	webhook := assertOutgoingWebhookExists(t, uid, *create)

	// Retrieve webhook
	webhooks, err := testDB.GetOutgoingWebhooksByUser(uid)
	assert.Nil(t, err)
	assert.True(t, len(webhooks) > 0)

	// Delete webhook
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	assert.Nil(t, err)

	// Unable to retrieve deleted webhook
	webhook, err = testDB.GetOutgoingWebhookByUserAndAlias(uid, &create.Alias)
	assert.Nil(t, err)
	assert.Nil(t, webhook)
}

func TestUpdateDefaultOutgoingWebhook(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create first default webhook
	uid := *testUser.ID
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	create := builder.Alias(
		"My default outgoing webhook",
	).Dummy().IsDefault(true).Build()
	first := assertOutgoingWebhookExists(t, uid, *create)

	// Create second default webhook
	builder = model.NewOutgoingWebhookCreateFormBuilder()
	create = builder.Alias(
		"My second outgoing webhook",
	).Dummy().IsDefault(true).Build()
	second := assertOutgoingWebhookExists(t, uid, *create)

	// Refresh first webhook
	first, err := testDB.GetOutgoingWebhookByID(*first.ID)
	assert.Nil(t, err)
	assert.False(t, first.IsDefault, "outgoing webhook should not be the default anymore")

	// Test default webhook query
	defaultOutgoingWebhook, err := testDB.GetOutgoingWebhookByUserAndAlias(uid, nil)
	assert.Nil(t, err)
	assert.NotNil(t, defaultOutgoingWebhook)
	assert.NotEqual(t, *first.ID, *defaultOutgoingWebhook.ID)
	assert.Equal(t, *second.ID, *defaultOutgoingWebhook.ID)

	// Cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *first.ID)
	assert.Nil(t, err)
	err = testDB.DeleteOutgoingWebhookByUser(uid, *second.ID)
	assert.Nil(t, err)
}
