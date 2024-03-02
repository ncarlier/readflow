package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
)

func TestCreateOutgoingWebhook(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new webhook
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	form := builder.Alias("test").Dummy().Build()

	webhook, err := service.Lookup().CreateOutgoingWebhook(testContext, *form)
	require.Nil(t, err)
	require.Equal(t, "test", webhook.Alias)
	require.NotEmpty(t, webhook.Config)

	// Create same webhook again
	_, err = service.Lookup().CreateOutgoingWebhook(testContext, *form)
	require.Equal(t, "already exists", err.Error())
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
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
			require.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
