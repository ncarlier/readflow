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

		article, err := reg.db.CreateOrUpdateArticle(*article)
		if err != nil {
			result.Errors = append(result.Errors, err)
		} else {
			result.Articles = append(result.Articles, article)
		}
	}
	var err error
	if len(result.Errors) > 0 {
		err = fmt.Errorf("Errors when creating articles")
	}
	return &result, err
}
