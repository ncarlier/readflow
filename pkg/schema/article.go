package schema

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/ncarlier/readflow/pkg/tooling"
)

var sortOrder = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "sortOrder",
		Description: "Sorting order",
		Values: graphql.EnumValueConfigMap{
			"asc": &graphql.EnumValueConfig{
				Value:       "asc",
				Description: "from older to newer",
			},
			"desc": &graphql.EnumValueConfig{
				Value:       "desc",
				Description: "from newer to older",
			},
		},
	},
)

var articleStatus = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "status",
		Description: "Article status",
		Values: graphql.EnumValueConfigMap{
			"read": &graphql.EnumValueConfig{
				Value:       "read",
				Description: "article is read",
			},
			"unread": &graphql.EnumValueConfig{
				Value:       "unread",
				Description: "article is not read",
			},
		},
	},
)

var articleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"html": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: articleStatus,
			},
			"starred": &graphql.Field{
				Type: graphql.Boolean,
			},
			"category": &graphql.Field{
				Type: categoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					article, ok := p.Source.(*model.Article)
					if !ok {
						return nil, errors.New("no article received by category resolver")
					}
					if article.CategoryID != nil {
						return service.Lookup().GetCategory(p.Context, *article.CategoryID)
					}
					return nil, nil
				},
			},
			"published_at": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var articlesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Articles",
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type: graphql.Int,
			},
			"endCursor": &graphql.Field{
				Type: graphql.Int,
			},
			"hasNext": &graphql.Field{
				Type: graphql.Boolean,
			},
			"entries": &graphql.Field{
				Type: graphql.NewList(articleType),
			},
		},
	},
)

var articleUpdateResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleUpdateResponseType",
		Fields: graphql.Fields{
			"article": &graphql.Field{
				Type: articleType,
			},
			"_all": &graphql.Field{
				Type: categoryType,
			},
		},
	},
)

// QUERIES

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
		"sortOrder": &graphql.ArgumentConfig{
			Description:  "sorting order of the entries",
			Type:         sortOrder,
			DefaultValue: "asc",
		},
	},
	Resolve: articlesResolver,
}

func articlesResolver(p graphql.ResolveParams) (interface{}, error) {
	pageRequest := model.ArticlesPageRequest{
		Limit:       tooling.GetGQLUintParameter("limit", p.Args),
		SortOrder:   tooling.GetGQLStringParameter("sortOrder", p.Args),
		AfterCursor: tooling.GetGQLUintParameter("afterCursor", p.Args),
		Category:    tooling.GetGQLUintParameter("category", p.Args),
		Status:      tooling.GetGQLStringParameter("status", p.Args),
		Starred:     tooling.GetGQLBoolParameter("starred", p.Args),
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
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}

	return service.Lookup().GetArticle(p.Context, id)
}

// MUTATIONS

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
		"starred": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
	},
	Resolve: updateArticleResolver,
}

func updateArticleResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}

	form := model.ArticleUpdateForm{
		ID:      id,
		Status:  tooling.GetGQLStringParameter("status", p.Args),
		Starred: tooling.GetGQLBoolParameter("starred", p.Args),
	}

	article, err := service.Lookup().UpdateArticle(p.Context, form)
	if err != nil {
		return nil, err
	}
	response := &model.ArticleStatusResponse{
		Article: article,
		All: &model.Category{
			Title: "_all",
		},
	}
	return response, nil
}

var markAllArticlesAsReadMutationField = &graphql.Field{
	Type:        graphql.NewList(categoryType),
	Description: "set all articles (of a category if provided) to read status",
	Args: graphql.FieldConfigArgument{
		"category": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: markAllArticlesAsReadResolver,
}

func markAllArticlesAsReadResolver(p graphql.ResolveParams) (interface{}, error) {
	var category *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["category"]); ok {
		category = &val
	}

	_, err := service.Lookup().MarkAllArticlesAsRead(p.Context, category)
	if err != nil {
		return nil, err
	}
	return categoriesResolver(p)
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
	var category *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["category"]); ok {
		category = &val
	}
	url, _ := p.Args["url"].(string)
	form := model.ArticleCreateForm{
		URL:        &url,
		CategoryID: category,
	}

	return service.Lookup().CreateArticle(p.Context, form, service.ArticleCreationOptions{})
}

var cleanHistoryMutationField = &graphql.Field{
	Type:        graphql.NewList(categoryType),
	Description: "remove all read articles",
	Resolve:     cleanHistoryResolver,
}

func cleanHistoryResolver(p graphql.ResolveParams) (interface{}, error) {
	_, err := service.Lookup().CleanHistory(p.Context)
	if err != nil {
		return nil, err
	}
	return categoriesResolver(p)
}
