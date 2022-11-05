package service

import (
	"context"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/scripting"
)

func mapArticleCreateFromToScriptInput(article *model.ArticleCreateForm) *scripting.ScriptInput {
	input := scripting.ScriptInput{
		Title: article.Title,
	}
	if article.URL != nil {
		input.URL = *article.URL
	}
	if article.Text != nil {
		input.Text = *article.Text
	}
	if article.HTML != nil {
		input.HTML = *article.HTML
	}
	if article.Tags != nil {
		input.Tags = strings.Split(*article.Tags, ",")
	}
	return &input
}

// ProcessArticleByScriptEngine apply user's script on the article
func (reg *Registry) ProcessArticleByScriptEngine(ctx context.Context, alias string, article *model.ArticleCreateForm) (scripting.OperationStack, error) {
	uid := getCurrentUserIDFromContext(ctx)

	noops := scripting.OperationStack{}

	// retrieve webhook
	webhook, err := reg.db.GetIncomingWebhookByUserAndAlias(uid, alias)
	if err != nil || webhook == nil {
		return noops, err
	}

	// limit execution time to 1 sec
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second))
	defer cancel()

	// build input script object
	input := mapArticleCreateFromToScriptInput(article)

	// exec user's script
	ops, err := reg.scriptEngine.Exec(ctx, webhook.Script, *input)
	if err != nil {
		return noops, err
	}

	// Apply setter
	for _, op := range ops {
		switch op.Name {
		case scripting.OpSetCategory:
			// Set category
			if cat, err := reg.db.GetCategoryByUserAndTitle(uid, op.Args[0]); err == nil {
				article.CategoryID = cat.ID
			}
		case scripting.OpSetTitle:
			// Set title
			article.Title = op.Args[0]
		}
	}
	return ops, err
}
