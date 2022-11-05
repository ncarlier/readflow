package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticleWithScriptEngine(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	cat := assertNewCategory(t)

	// Create incoming webhook
	script := fmt.Sprintf(`
if ( Title ~= /script/i ) {
	setCategory("%s");
	return true;
}
return false;
`, cat.Title)
	assertNewIncomingWebhook(t, "foo", script)
	ctx := context.WithValue(testContext, constant.ContextIncomingWebhookAlias, "foo")

	// Create article
	form := model.ArticleCreateForm{
		Title: "TestCreateArticleWithScriptEngine",
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(ctx, form, opts)
	assert.Nil(t, err)
	assert.Equal(t, form.Title, art.Title)
	assert.Equal(t, *testUser.ID, art.UserID)
	assert.NotNil(t, art.CategoryID)
	assert.Equal(t, *cat.ID, *art.CategoryID)
}
