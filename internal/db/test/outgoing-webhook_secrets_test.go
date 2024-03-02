package dbtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/secret"
)

func TestOutgoingWebhookSecretsUpdate(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// create webhook
	uid := *testUser.ID
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	create := builder.Alias(
		"My test outgoing webhook with secrets",
	).Dummy().Build()
	webhook := assertOutgoingWebhookExists(t, uid, *create)

	require.Equal(t, "bar", webhook.Secrets["foo"])

	// update webhook with empty secret should not affect secret value
	secrets := make(secret.Secrets)
	secrets["foo"] = ""
	update := model.OutgoingWebhookUpdateForm{
		ID:      *webhook.ID,
		Secrets: &secrets,
	}
	webhook, err := testDB.UpdateOutgoingWebhookForUser(uid, update)
	require.Nil(t, err)
	require.NotNil(t, webhook)
	require.Equal(t, "bar", webhook.Secrets["foo"])

	// update webhook secrets by omitting a previous one should remove it
	delete(*update.Secrets, "foo")
	(*update.Secrets)["zoo"] = "baz"
	webhook, err = testDB.UpdateOutgoingWebhookForUser(uid, update)
	require.Nil(t, err)
	require.NotNil(t, webhook)
	require.Equal(t, "baz", webhook.Secrets["zoo"])
	require.NotContains(t, "bar", webhook.Secrets)

	// cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	require.Nil(t, err)
}

func TestOutgoingWebhookSecretsManagement(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// create secret engine
	engine, err := secret.NewSecretsEngineProvider("file://./assets/secret.key")
	require.Nil(t, err)
	require.NotNil(t, engine)

	// create webhook
	uid := *testUser.ID
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	create := builder.Alias(
		"My test outgoing webhook with un managed secrets",
	).Dummy().Build()
	webhook := assertOutgoingWebhookExists(t, uid, *create)

	require.Equal(t, "bar", webhook.Secrets["foo"])

	// seal all secrets
	ctx := context.Background()
	nb, err := testDB.ManageOutgoingWebhookSecrets(ctx, engine, secret.Seal)
	require.Nil(t, err)
	require.Greater(t, nb, uint(0))

	// refresh webhook
	webhook, err = testDB.GetOutgoingWebhookByID(*webhook.ID)
	assert.Nil(t, err)
	assert.NotEqual(t, "bar", webhook.Secrets["foo"])

	// un-seal all secrets
	ctx = context.Background()
	nb, err = testDB.ManageOutgoingWebhookSecrets(ctx, engine, secret.UnSeal)
	require.Nil(t, err)
	require.Greater(t, nb, uint(0))

	// refresh webhook
	webhook, err = testDB.GetOutgoingWebhookByID(*webhook.ID)
	assert.Nil(t, err)
	assert.Equal(t, "bar", webhook.Secrets["foo"])

	// cleanup
	err = testDB.DeleteOutgoingWebhookByUser(uid, *webhook.ID)
	require.Nil(t, err)
}
