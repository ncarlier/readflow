package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestCreateOutgoingWebhook(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new webhook
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	form := builder.Alias("test").Dummy().Build()

	webhook, err := service.Lookup().CreateOutgoingWebhook(testContext, *form)
	assert.Nil(t, err)
	assert.Equal(t, "test", webhook.Alias)
	assert.NotEmpty(t, webhook.Config)

	// Create same webhook again
	_, err = service.Lookup().CreateOutgoingWebhook(testContext, *form)
	assert.Equal(t, "already exists", err.Error())
}

func TestCreateoutgoingWebhooksExceedingQuota(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create 3 webhooks (quota is 2)
	for i := 1; i <= 3; i++ {
		builder := model.NewOutgoingWebhookCreateFormBuilder()
		form := builder.Alias(fmt.Sprintf("TestCreateWebhooksExceedingQuota-%d", i)).Dummy().Build()
		_, err := service.Lookup().CreateOutgoingWebhook(testContext, *form)
		if i <= 2 {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
