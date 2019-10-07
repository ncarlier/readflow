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

var articleStatusResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleStatusResponseType",
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
		"sortOrder": &graphql.ArgumentConfig{
			Description:  "sorting order of the entries",
			Type:         sortOrder,
			DefaultValue: "asc",
		},
	},
	Resolve: articlesResolver,
}

func articlesResolver(p graphql.ResolveParams) (interface{}, error) {
	sortOrder, _ := p.Args["sortOrder"].(string)
	var limit uint
	if val, ok := tooling.ConvGQLIntToUint(p.Args["limit"]); ok {
		limit = val
	}
	var category *uint
	if val, ok := tooling.ConvGQLIntToUint(p.Args["category"]); ok {
		category = &val
	}
	var afterCursor *uint
	if val, ok := tooling.ConvGQLIntToUint(p.Args["afterCursor"]); ok {
		afterCursor = &val
	}
	var status *string
	if val, ok := p.Args["status"].(string); ok {
		status = &val
	}

	pageRequest := model.ArticlesPageRequest{
		Limit:       limit,
		SortOrder:   sortOrder,
		AfterCursor: afterCursor,
		Category:    category,
		Status:      status,
	}

	articles, err := service.Lookup().GetArticles(p.Context, pageRequest)
	if err != nil {
		return nil, err
	}
	return articles, nil
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

	article, err := service.Lookup().GetArticle(p.Context, id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

// MUTATIONS

var updateArticleStatusMutationField = &graphql.Field{
	Type:        articleStatusResponseType,
	Description: "update article status (read or unread)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"status": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(articleStatus),
		},
	},
	Resolve: updateArticleStatusResolver,
}

func updateArticleStatusResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}
	status, _ := p.Args["status"].(string)

	article, err := service.Lookup().UpdateArticleStatus(p.Context, id, status)
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
	form := model.ArticleForm{
		URL:        &url,
		CategoryID: category,
	}

	article, err := service.Lookup().CreateArticle(p.Context, form, service.ArticleCreationOptions{})
	if err != nil {
		return nil, err
	}
	return article, nil
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
