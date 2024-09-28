package article

import (
	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/schema/category"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
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
			Type: graphql.ID,
		},
		"status": &graphql.ArgumentConfig{
			Type: articleStatus,
		},
		"stars": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"refresh": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
	},
	Resolve: updateArticleResolver,
}

func updateArticleResolver(p graphql.ResolveParams) (interface{}, error) {
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}

	form := model.ArticleUpdateForm{
		ID:             *id,
		Title:          utils.ParseGraphQLArgument[string](p.Args, "title"),
		Text:           utils.ParseGraphQLArgument[string](p.Args, "text"),
		CategoryID:     utils.ParseGraphQLID(p.Args, "category_id"),
		Status:         utils.ParseGraphQLArgument[string](p.Args, "status"),
		Stars:          utils.ParseGraphQLArgument[int](p.Args, "stars"),
		RefreshContent: utils.ParseGraphQLArgument[bool](p.Args, "refresh"),
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
	categoryID := utils.ParseGraphQLID(p.Args, "category")
	status := utils.ParseGraphQLArgument[string](p.Args, "status")
	if status == nil || *status == "read" {
		return nil, global.InvalidParameterError("status")
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
	form := model.ArticleCreateForm{
		URL:        utils.ParseGraphQLArgument[string](p.Args, "url"),
		CategoryID: utils.ParseGraphQLID(p.Args, "category"),
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
