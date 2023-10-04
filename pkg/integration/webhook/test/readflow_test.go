package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/all"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/secret"
	"github.com/stretchr/testify/require"
)

func TestReadflowWebhook(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, constant.UserAgent, r.Header.Get("User-Agent"))
		require.Equal(t, constant.ContentTypeJSON, r.Header.Get("Content-Type"))
		require.Equal(t, "/articles", r.URL.Path)
		username, password, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "api", username)
		require.Equal(t, "foo", password)
	}))
	defer srv.Close()

	secrets := make(secret.Secrets)
	secrets["api_key"] = "foo"

	outgoingWebhook := model.OutgoingWebhook{
		Provider: "readflow",
		Config:   fmt.Sprintf(`{"endpoint": %q}`, srv.URL),
		Secrets:  secrets,
	}
	conf := config.NewConfig()
	conf.HTTP.PublicURL = "http://localhost:3000"

	provider, err := webhook.NewOutgoingWebhookProvider(outgoingWebhook, *conf)
	require.Nil(t, err)
	require.NotNil(t, provider)

	result, err := provider.Send(context.TODO(), article)
	require.Nil(t, err)
	require.NotNil(t, result)
}
