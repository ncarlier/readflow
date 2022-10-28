package article

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/schema/category"
	"github.com/ncarlier/readflow/pkg/service"
)

var updateArticleMutationField = &graphql.Field{
	Type:        articleUpdateResponseType,
	Description: "update article",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"text": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"category_id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"status": &graphql.ArgumentConfig{
			Type: articleStatus,
		},
		"stars": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: updateArticleResolver,
}

func updateArticleResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}

	form := model.ArticleUpdateForm{
		ID:         id,
		Title:      helper.GetGQLStringParameter("title", p.Args),
		Text:       helper.GetGQLStringParameter("text", p.Args),
		CategoryID: helper.GetGQLUintParameter("category_id", p.Args),
		Status:     helper.GetGQLStringParameter("status", p.Args),
		Stars:      helper.GetGQLUintParameter("stars", p.Args),
	}

	article, err := service.Lookup().UpdateArticle(p.Context, form)
	if err != nil {
		return nil, err
	}
	return struct {
		Article *model.Article
	}{
		Article: article,
	}, nil
}

var markAllArticlesAsReadMutationField = &graphql.Field{
	Type:        category.ListResponseType,
	Description: "set all articles (of a category if provided) to read status",
	Args: graphql.FieldConfigArgument{
		"status": &graphql.ArgumentConfig{
			Type: articleStatus,
		},
		"category": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: markAllArticlesAsReadResolver,
}

func markAllArticlesAsReadResolver(p graphql.ResolveParams) (interface{}, error) {
	var categoryID *uint
	if val, ok := helper.ConvGQLStringToUint(p.Args["category"]); ok {
		categoryID = &val
	}
	status := helper.GetGQLStringParameter("status", p.Args)
	if status == nil || *status == "read" {
		return nil, errors.New("invalid status")
	}

	_, err := service.Lookup().MarkAllArticlesAsRead(p.Context, *status, categoryID)
	if err != nil {
		return nil, err
	}
	return category.CategoriesResolver(p)
}

var addArticleMutationField = &graphql.Field{
	Type:        articleType,
	Description: "add new article",
	Args: graphql.FieldConfigArgument{
		"url": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"category": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: addArticleResolver,
}

func addArticleResolver(p graphql.ResolveParams) (interface{}, error) {
	var categoryID *uint
	if val, ok := helper.ConvGQLStringToUint(p.Args["category"]); ok {
		categoryID = &val
	}
	url, _ := p.Args["url"].(string)
	form := model.ArticleCreateForm{
		URL:        &url,
		CategoryID: categoryID,
	}

	return service.Lookup().CreateArticle(p.Context, form, service.ArticleCreationOptions{})
}

var cleanHistoryMutationField = &graphql.Field{
	Type:        category.ListResponseType,
	Description: "remove all read articles",
	Resolve:     cleanHistoryResolver,
}

func cleanHistoryResolver(p graphql.ResolveParams) (interface{}, error) {
	_, err := service.Lookup().CleanHistory(p.Context)
	if err != nil {
		return nil, err
	}
	return category.CategoriesResolver(p)
}

func init() {
	schema.AddMutationField("addArticle", addArticleMutationField)
	schema.AddMutationField("updateArticle", updateArticleMutationField)
	schema.AddMutationField("markAllArticlesAsRead", markAllArticlesAsReadMutationField)
	schema.AddMutationField("cleanHistory", cleanHistoryMutationField)
}
