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
		ID:     id,
		Status: helper.GetGQLStringParameter("status", p.Args),
		Stars:  helper.GetGQLUintParameter("stars", p.Args),
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

	_, err := service.Lookup().MarkAllArticlesAsRead(p.Context, categoryID)
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
	Type:        graphql.NewList(category.Type),
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
