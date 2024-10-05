package article

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
)

var articlesQueryField = &graphql.Field{
	Type: articlesType,
	Args: graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Description:  "max number of entries to returns",
			Type:         graphql.Int,
			DefaultValue: 10,
		},
		"afterCursor": &graphql.ArgumentConfig{
			Description: "retrive entries after this cursor",
			Type:        graphql.Int,
		},
		"category": &graphql.ArgumentConfig{
			Description: "filter entries by this category",
			Type:        graphql.Int,
		},
		"status": &graphql.ArgumentConfig{
			Description: "filter entries by this status",
			Type:        articleStatus,
		},
		"starred": &graphql.ArgumentConfig{
			Description: "filter entries by this starred value",
			Type:        graphql.Boolean,
		},
		"query": &graphql.ArgumentConfig{
			Description: "filter entries by full-text search",
			Type:        graphql.String,
		},
		"sortOrder": &graphql.ArgumentConfig{
			Description:  "sorting order of the entries",
			Type:         sortOrder,
			DefaultValue: "asc",
		},
		"sortBy": &graphql.ArgumentConfig{
			Description:  "sorting attribute of the entries",
			Type:         sortBy,
			DefaultValue: "key",
		},
	},
	Resolve: articlesResolver,
}

func articlesResolver(p graphql.ResolveParams) (interface{}, error) {
	pageRequest := model.ArticlesPageRequest{
		Limit:       utils.ParseGraphQLArgument[int](p.Args, "limit"),
		SortOrder:   utils.ParseGraphQLArgument[string](p.Args, "sortOrder"),
		SortBy:      utils.ParseGraphQLArgument[string](p.Args, "sortBy"),
		AfterCursor: utils.ParseGraphQLID(p.Args, "afterCursor"),
		Category:    utils.ParseGraphQLID(p.Args, "category"),
		Status:      utils.ParseGraphQLArgument[string](p.Args, "status"),
		Starred:     utils.ParseGraphQLArgument[bool](p.Args, "starred"),
		Query:       utils.ParseGraphQLArgument[string](p.Args, "query"),
	}

	return service.Lookup().GetArticles(p.Context, pageRequest)
}

var articleQueryField = &graphql.Field{
	Type: articleType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: articleResolver,
}

func articleResolver(p graphql.ResolveParams) (interface{}, error) {
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}

	return service.Lookup().GetArticle(p.Context, *id)
}

func thumbnailsResolver(p graphql.ResolveParams) (interface{}, error) {
	article, ok := p.Source.(*model.Article)
	if !ok {
		return nil, errors.New("thumbnails resolver is expecting an article")
	}
	return service.Lookup().GetArticleThumbnailHashSet(article), nil
}

func init() {
	schema.AddQueryField("articles", articlesQueryField)
	schema.AddQueryField("article", articleQueryField)
}
