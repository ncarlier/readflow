package schema

import (
	"errors"

	"github.com/ncarlier/reader/pkg/model"

	"github.com/ncarlier/reader/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/service"
)

var ruleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Rule",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"category_id": &graphql.Field{
				Type: graphql.Int,
			},
			"rule": &graphql.Field{
				Type: graphql.String,
			},
			"priority": &graphql.Field{
				Type: graphql.Int,
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

// QUERIES

var rulesQueryField = &graphql.Field{
	Type:    graphql.NewList(ruleType),
	Resolve: rulesResolver,
}

func rulesResolver(p graphql.ResolveParams) (interface{}, error) {
	rules, err := service.Lookup().GetRules(p.Context)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

var ruleQueryField = &graphql.Field{
	Type: ruleType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: ruleResolver,
}

func ruleResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid rule ID")
	}
	rule, err := service.Lookup().GetRule(p.Context, id)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// MUTATIONS

var createOrUpdateRuleMutationField = &graphql.Field{
	Type:        ruleType,
	Description: "create or update a rule (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"alias": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"rule": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"priority": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 1,
		},
		"category_id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: createOrUpdateRuleResolver,
}

func createOrUpdateRuleResolver(p graphql.ResolveParams) (interface{}, error) {
	var id *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["id"]); ok {
		id = &val
	}
	categoryID, ok := tooling.ConvGQLIntToUint(p.Args["category_id"])
	if !ok {
		return nil, errors.New("invalid category ID")
	}

	alias, _ := p.Args["alias"].(string)
	ruleDef, _ := p.Args["rule"].(string)
	priority, _ := p.Args["priority"].(int)

	form := model.RuleForm{
		ID:         id,
		Alias:      alias,
		CategoryID: categoryID,
		Rule:       ruleDef,
		Priority:   priority,
	}

	rule, err := service.Lookup().CreateOrUpdateRule(p.Context, form)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

var deleteRulesMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete rules",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteRulesResolver,
}

func deleteRulesResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid rule ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := tooling.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	nb, err := service.Lookup().DeleteRules(p.Context, ids)
	if err != nil {
		return nil, err
	}
	return nb, nil
}
