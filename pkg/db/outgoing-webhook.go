package db

import (
	"context"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/secret"
)

// OutgoingWebhookRepository is the repository interface to manage outgoing webhooks
type OutgoingWebhookRepository interface {
	GetOutgoingWebhookByID(id uint) (*model.OutgoingWebhook, error)
	GetOutgoingWebhookByUserAndAlias(uid uint, alias *string) (*model.OutgoingWebhook, error)
	GetOutgoingWebhooksByUser(uid uint) ([]model.OutgoingWebhook, error)
	CreateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookCreateForm) (*model.OutgoingWebhook, error)
	UpdateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookUpdateForm) (*model.OutgoingWebhook, error)
	DeleteOutgoingWebhookByUser(uid uint, id uint) error
	DeleteOutgoingWebhooksByUser(uid uint, ids []uint) (int64, error)
	ManageOutgoingWebhookSecrets(ctx context.Context, provider secret.EngineProvider, action secret.Action) (uint, error)
}
