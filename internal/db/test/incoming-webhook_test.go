package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/internal/model"
)

func assertIncomingWebhookExists(t *testing.T, uid uint, alias, script string) *model.IncomingWebhook {
	webhook, err := testDB.GetIncomingWebhookByUserAndAlias(uid, alias)
	assert.Nil(t, err)
	if webhook != nil {
		return webhook
	}

	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias(alias).Script(script).Build()

	webhook, err = testDB.CreateIncomingWebhookForUser(uid, *form)
	assert.Nil(t, err)
	assert.NotNil(t, webhook)
	assert.NotNil(t, webhook.ID)
	assert.Equal(t, alias, webhook.Alias)
	assert.Equal(t, script, webhook.Script)
	assert.NotEqual(t, "", webhook.Token)
	return webhook
}
func TestCreateOrUpdateIncomingWebhook(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My test incoming webhook"
	script := "return false;"

	assertIncomingWebhookExists(t, uid, alias, script)
}

func TestDeleteIncomingWebhook(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID
	alias := "My incoming webhook"
	script := "return false;"

	// Assert webhook exists
	webhook := assertIncomingWebhookExists(t, uid, alias, script)

	err := testDB.DeleteIncomingWebhookByUser(uid, *webhook.ID)
	assert.Nil(t, err)

	webhook, err = testDB.GetIncomingWebhookByToken(webhook.Token)
	assert.Nil(t, err)
	assert.Nil(t, webhook)
}

func TestGetIncomingWebhooksByUserID(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	uid := *testUser.ID

	webhooks, err := testDB.GetIncomingWebhooksByUser(uid)
	assert.Nil(t, err)
	assert.NotNil(t, webhooks)
	assert.Positive(t, len(webhooks))
}
