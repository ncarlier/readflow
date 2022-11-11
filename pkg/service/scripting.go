package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/helper"
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
	if article.Origin != nil {
		input.Origin = *article.Origin
	}
	if article.Tags != nil {
		input.Tags = strings.Split(*article.Tags, ",")
	}
	return &input
}

// processArticleByScriptEngine apply user's script on the article
func (reg *Registry) processArticleByScriptEngine(ctx context.Context, alias string, article *model.ArticleCreateForm) (scripting.OperationStack, error) {
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

	return ops, err
}

func (reg *Registry) execSetOperations(ctx context.Context, ops scripting.OperationStack, article *model.ArticleCreateForm) {
	uid := getCurrentUserIDFromContext(ctx)
	for _, op := range ops {
		value := op.Args[0]
		switch op.Name {
		case scripting.OpSetCategory:
			// set category
			if cat, err := reg.db.GetCategoryByUserAndTitle(uid, value); err == nil && cat != nil {
				article.CategoryID = cat.ID
			}
		case scripting.OpSetText:
			// set text
			article.Text = &value
		case scripting.OpSetTitle:
			// set title
			article.Title = value
		}
	}
}

func (reg *Registry) execOtherOperations(ctx context.Context, ops scripting.OperationStack, article *model.Article) error {
	for _, op := range ops {
		switch op.Name {
		case scripting.OpSendNotification:
			// build notification
			href := fmt.Sprintf("/inbox/%d", article.ID)
			if article.CategoryID != nil {
				href = fmt.Sprintf("/categories/%d/%d", *article.CategoryID, article.ID)
			}
			notif := &model.DeviceNotification{
				Title: "New article to read",
				Body:  helper.Truncate(article.Title, 64),
				Href:  href,
			}
			if _, err := reg.NotifyDevices(ctx, notif); err != nil {
				return err
			}
		case scripting.OpTriggerWebhook:
			name := op.Args[0]
			if err := reg.SendArticle(ctx, article.ID, &name); err != nil {
				return err
			}
		}
	}
	return nil
}
