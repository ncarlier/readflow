package service

import (
	"context"
	"fmt"

	"github.com/ncarlier/reader/pkg/constant"
	"github.com/ncarlier/reader/pkg/model"
)

// CreateArticles creates new articles
func (reg *Registry) CreateArticles(ctx context.Context, data []model.ArticleForm) (*model.Articles, error) {
	userID := ctx.Value(constant.UserID).(uint32)
	result := model.Articles{}
	for _, art := range data {
		builder := model.NewArticleBuilder()
		article := builder.UserID(
			userID,
		).Form(&art).Build()

		// TODO validate article!

		reg.logger.Debug().Uint32(
			"uid", userID,
		).Str("title", article.Title).Msg("creating article...")
		article, err := reg.db.CreateOrUpdateArticle(*article)
		if err != nil {
			result.Errors = append(result.Errors, err)
			reg.logger.Info().Err(err).Uint32(
				"uid", userID,
			).Str("title", art.Title).Msg("unable to create article")
		} else {
			result.Articles = append(result.Articles, article)
			reg.logger.Debug().Uint32(
				"uid", userID,
			).Str("title", article.Title).Uint32("id", *article.ID).Msg("article created")
		}
	}
	var err error
	if len(result.Errors) > 0 {
		err = fmt.Errorf("Errors when creating articles")
	}
	return &result, err
}
