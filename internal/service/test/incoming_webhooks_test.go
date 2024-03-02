package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
)

func TestCreateIncomingWebhook(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new webhook
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias("test").Build()

	webhook, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
	require.Nil(t, err)
	require.Equal(t, "test", webhook.Alias)
	require.NotEmpty(t, webhook.Token)

	// Create same webhook again
	_, err = service.Lookup().CreateIncomingWebhook(testContext, *form)
	require.Equal(t, "already exists", err.Error())
}

func TestCreateIncomingWebhooksExceedingQuota(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create 3 webhooks (quota is 2)
	for i := 1; i <= 3; i++ {
		builder := model.NewIncomingWebhookCreateFormBuilder()
		form := builder.Alias(fmt.Sprintf("TestCreateWebhooksExceedingQuota-%d", i)).Build()
		_, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
		if i <= 2 {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
			require.Equal(t, service.ErrUserQuotaReached.Error(), err.Error())
		}
	}
}
