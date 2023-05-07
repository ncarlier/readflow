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
	"github.com/stretchr/testify/require"
)

var article = model.Article{
	ID:    uint(1),
	Title: "Foo & Bar",
}

var body = "title={{ title | urlquery }}"

var headers = fmt.Sprintf(`{"Content-Type": "%s", "X-API-Key": "xxx"}`, constant.ContentTypeForm)

func TestGenericWebhook(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, constant.UserAgent, r.Header.Get("User-Agent"))
		require.Equal(t, constant.ContentTypeForm, r.Header.Get("Content-Type"))
		require.Equal(t, "xxx", r.Header.Get("X-API-Key"))
		require.Equal(t, article.Title, r.FormValue("title"))
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", "https://example.org")
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	outgoingWebhook := model.OutgoingWebhook{
		Provider: "generic",
		Config:   fmt.Sprintf(`{"endpoint": "%s", "headers": %s, "body": "%s"}`, srv.URL, headers, body),
	}
	conf := config.NewConfig()
	conf.Global.PublicURL = "http://localhost:3000"

	provider, err := webhook.NewOutgoingWebhookProvider(outgoingWebhook, *conf)
	require.Nil(t, err)
	require.NotNil(t, provider)

	result, err := provider.Send(context.TODO(), article)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Nil(t, result.JSON)
	require.NotNil(t, result.Text)
	require.Equal(t, "ok", *result.Text)
	require.NotNil(t, result.URL)
	require.Equal(t, "https://example.org", *result.URL)
}
