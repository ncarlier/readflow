package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/stretchr/testify/assert"
)

func assertNewArticle(t *testing.T, title string) *model.Article {
	req := model.ArticleCreateForm{
		Title: title,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	assert.Nil(t, err)
	assert.Equal(t, title, art.Title)
	assert.Equal(t, *testUser.ID, art.UserID)

	return art
}

func assertNewCategory(t *testing.T) *model.Category {
	formBuilder := model.NewCategoryCreateFormBuilder()
	form := formBuilder.Random().Build()
	cat, err := service.Lookup().CreateCategory(testContext, *form)
	assert.Nil(t, err)
	assert.Equal(t, form.Title, cat.Title)
	assert.Equal(t, *testUser.ID, *cat.UserID)
	return cat
}

func assertNewIncomingWebhook(t *testing.T, alias, script string) *model.IncomingWebhook {
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias(alias).Script(script).Build()
	webhook, err := service.Lookup().CreateIncomingWebhook(testContext, *form)
	assert.Nil(t, err)
	assert.Equal(t, *testUser.ID, webhook.UserID)
	assert.Equal(t, alias, webhook.Alias)
	assert.Equal(t, script, webhook.Script)
	assert.NotEmpty(t, webhook.Token)
	return webhook
}
