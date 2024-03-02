package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCreateArticleWithScriptEngine(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	cat := requireNewCategory(t)

	// Create incoming webhook
	script := fmt.Sprintf(`
if ( Title ~= /script/i ) {
	setCategory("%s");
	return true;
}
return false;
`, cat.Title)
	webhook := requireNewIncomingWebhook(t, "foo", script)
	ctx := context.WithValue(testContext, global.ContextIncomingWebhook, webhook)

	// Create article
	form := model.ArticleCreateForm{
		Title: "TestCreateArticleWithScriptEngine",
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(ctx, form, opts)
	require.Nil(t, err)
	require.Equal(t, form.Title, art.Title)
	require.Equal(t, *testUser.ID, art.UserID)
	require.NotNil(t, art.CategoryID)
	require.Equal(t, *cat.ID, *art.CategoryID)
}
