package test

import (
	"testing"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/stretchr/testify/require"
)

func requireNewArticle(t *testing.T, title string) *model.Article {
	req := model.ArticleCreateForm{
		Title: title,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	require.Nil(t, err)
	require.Equal(t, title, art.Title)
	require.Equal(t, *testUser.ID, art.UserID)

	return art
}

func requireNewCategory(t *testing.T) *model.Category {
	formBuilder := model.NewCategoryCreateFormBuilder()
	form := formBuilder.Random().Build()
	cat, err := service.Lookup().CreateCategory(testContext, *form)
	require.Nil(t, err)
	require.Equal(t, form.Title, cat.Title)
	require.Equal(t, *testUser.ID, *cat.UserID)
	return cat
}

func requireNewIncomingWebhook(t *testing.T, alias, script string) *model.IncomingWebhook {
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias(alias).Script(script).Build()
	webhook, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
	require.Nil(t, err)
	require.Equal(t, *testUser.ID, webhook.UserID)
	require.Equal(t, alias, webhook.Alias)
	require.Equal(t, script, webhook.Script)
	require.NotEmpty(t, webhook.Token)
	return webhook
}
