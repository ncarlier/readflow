package incomingwebhook

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
)

var incomingWebhookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "IncomingWebhook",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					webhook, ok := p.Source.(*model.IncomingWebhook)
					if !ok {
						return nil, errors.New("unsuported type received by email resolver")
					}
					hashid := service.Lookup().GetUserHashID(webhook.UserID)
					hostname := service.Lookup().GetConfig().SMTP.Hostname
					email := fmt.Sprintf("%s-%s@%s", webhook.Alias, hashid, hostname)
					if _, err := mail.ParseAddress(email); err != nil {
						return nil, nil
					}
					return email, nil
				},
			},
			"script": &graphql.Field{
				Type: graphql.String,
			},
			"last_usage_at": &graphql.Field{
				Type: graphql.DateTime,
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
