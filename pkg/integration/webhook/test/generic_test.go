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
	"github.com/stretchr/testify/assert"
)

var article = model.Article{
	ID:    uint(1),
	Title: "Foo & Bar",
}

var body = "title={{ title | urlquery }}"

var headers = fmt.Sprintf(`{"Content-Type": "%s", "X-API-Key": "xxx"}`, constant.ContentTypeForm)

func TestGenericWebhook(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, constant.UserAgent, r.Header.Get("User-Agent"))
		assert.Equal(t, constant.ContentTypeForm, r.Header.Get("Content-Type"))
		assert.Equal(t, "xxx", r.Header.Get("X-API-Key"))
		assert.Equal(t, article.Title, r.FormValue("title"))
	}))
	defer srv.Close()

	outgoingWebhook := model.OutgoingWebhook{
		Provider: "generic",
		Config:   fmt.Sprintf(`{"endpoint": "%s", "headers": %s, "body": "%s"}`, srv.URL, headers, body),
	}
	conf := config.NewConfig()
	conf.Global.PublicURL = "http://localhost:3000"

	provider, err := webhook.NewOutgoingWebhookProvider(outgoingWebhook, *conf)
	assert.Nil(t, err)
	assert.NotNil(t, provider)

	err = provider.Send(context.TODO(), article)
	assert.Nil(t, err)
}
