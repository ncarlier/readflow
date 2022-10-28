package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestCreateIncomingWebHook(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create new webhook
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias("test").Build()

	webhook, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
	assert.Nil(t, err)
	assert.Equal(t, "test", webhook.Alias)
	assert.NotEmpty(t, webhook.Token)

	// Create same webhook again
	_, err = service.Lookup().CreateIncomingWebhook(testContext, *form)
	assert.Equal(t, "already exists", err.Error())
}

func TestCreateIncomingWebhooksExceedingQuota(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create 3 webhooks (quota is 2)
	for i := 1; i <= 3; i++ {
		builder := model.NewIncomingWebhookCreateFormBuilder()
		form := builder.Alias(fmt.Sprintf("TestCreateWebhooksExceedingQuota-%d", i)).Build()
		_, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
		if i <= 2 {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
