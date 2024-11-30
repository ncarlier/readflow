package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/scripting"
	"github.com/ncarlier/readflow/pkg/utils"
)

func mapArticleCreateFromToScriptInput(article *model.ArticleCreateForm) *scripting.ScriptInput {
	input := scripting.ScriptInput{
		Title:  article.Title,
		URL:    utils.PtrValueOr(article.URL, ""),
		HTML:   utils.PtrValueOr(article.HTML, ""),
		Text:   utils.PtrValueOr(article.Text, ""),
		Origin: utils.PtrValueOr(article.Origin, ""),
	}
	if article.Tags != nil {
		input.Tags = strings.Split(*article.Tags, ",")
	}
	if hashtags := article.Hashtags(); len(hashtags) > 0 {
		input.Tags = append(input.Tags, hashtags...)
	}
	return &input
}

// processArticleByScriptEngine apply user's script on the article
func (reg *Registry) processArticleByScriptEngine(ctx context.Context, webhook *model.IncomingWebhook, article *model.ArticleCreateForm) (scripting.OperationStack, error) {
	noops := scripting.OperationStack{}

	if webhook == nil {
		return noops, nil
	}

	// limit execution time to 1 sec
	// TODO inherit context timeout
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
	category := ""
	for _, op := range ops {
		switch op.Name {
		case scripting.OpSetCategory:
			// only execute last setCategory operation
			category = op.GetFirstArg()
		case scripting.OpSetText:
			// set text
			value := op.GetFirstArg()
			article.Text = &value
		case scripting.OpSetHTML:
			// set HTML
			value := op.GetFirstArg()
			article.HTML = &value
		case scripting.OpSetTitle:
			// set title
			article.Title = op.GetFirstArg()
		}
	}
	if category != "" {
		if cat, err := reg.db.GetCategoryByUserAndTitle(uid, category); err == nil && cat != nil {
			article.CategoryID = cat.ID
		}
	}
}

func (reg *Registry) execOtherOperations(ctx context.Context, ops scripting.OperationStack, article *model.Article) error {
	// allows only 2 webhook trigger
	hardLimitCounter := 2
	// status to update if asked
	status := ""
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
				Body:  utils.Truncate(article.Title, 64),
				Href:  href,
			}
			if _, err := reg.NotifyDevices(ctx, notif); err != nil {
				return err
			}
		case scripting.OpTriggerWebhook:
			if hardLimitCounter == 0 {
				continue
			}
			hardLimitCounter--
			name := op.GetFirstArg()
			if _, err := reg.SendArticle(ctx, article.ID, &name); err != nil {
				return err
			}
		case scripting.OpMarkAsRead:
			status = "read"
		case scripting.OpMarkAsToRead:
			status = "to_read"
		}
	}
	if status != "" {
		// TODO move this logic to execSetOperations step
		update := model.ArticleUpdateForm{
			ID:     article.ID,
			Status: &status,
		}
		if _, err := reg.UpdateArticle(ctx, update); err != nil {
			return err
		}
	}
	return nil
}
