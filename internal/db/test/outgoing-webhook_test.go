package dbtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
)

func assertOutgoingWebhookExists(t *testing.T, uid uint, form model.OutgoingWebhookCreateForm) *model.OutgoingWebhook {
	result, err := testDB.GetOutgoingWebhookByUserAndAlias(uid, &form.Alias)
	require.Nil(t, err)
	if result != nil {
		return result
	}

	result, err = testDB.CreateOutgoingWebhookForUser(uid, form)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.ID)
	require.Equal(t, uid, *result.UserID)
	require.Equal(t, form.Alias, result.Alias)
	require.Equal(t, form.Provider, result.Provider)
	require.Equal(t, form.Config, result.Config)
	for k := range form.Secrets {
		require.Contains(t, result.Secrets, k)
	}
	require.Equal(t, form.IsDefault, result.IsDefault)
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
	require.Nil(t, err)
	require.NotNil(t, webhook)
	require.NotNil(t, webhook.ID)
	require.Equal(t, alias, webhook.Alias)
	require.Equal(t, create.Provider, webhook.Provider)
	require.Equal(t, create.Config, webhook.Config)
	require.Equal(t, create.IsDefault, webhook.IsDefault)

	// Cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	require.Nil(t, err)
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
	require.Nil(t, err)
	require.True(t, len(webhooks) > 0)

	// Delete webhook
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	require.Nil(t, err)

	// Unable to retrieve deleted webhook
	webhook, err = testDB.GetOutgoingWebhookByUserAndAlias(uid, &create.Alias)
	require.Nil(t, err)
	require.Nil(t, webhook)
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
	require.Nil(t, err)
	require.False(t, first.IsDefault, "outgoing webhook should not be the default anymore")

	// Test default webhook query
	defaultOutgoingWebhook, err := testDB.GetOutgoingWebhookByUserAndAlias(uid, nil)
	require.Nil(t, err)
	require.NotNil(t, defaultOutgoingWebhook)
	require.NotEqual(t, *first.ID, *defaultOutgoingWebhook.ID)
	require.Equal(t, *second.ID, *defaultOutgoingWebhook.ID)

	// Cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *first.ID)
	require.Nil(t, err)
	err = testDB.DeleteOutgoingWebhookByUser(uid, *second.ID)
	require.Nil(t, err)
}
