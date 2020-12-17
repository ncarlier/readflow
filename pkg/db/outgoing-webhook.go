package db

import "github.com/ncarlier/readflow/pkg/model"

// OutgoingWebhookRepository is the repository interface to manage outgoing webhooks
type OutgoingWebhookRepository interface {
	GetOutgoingWebhookByID(id uint) (*model.OutgoingWebhook, error)
	GetOutgoingWebhookByUserAndAlias(uid uint, alias *string) (*model.OutgoingWebhook, error)
	GetOutgoingWebhooksByUser(uid uint) ([]model.OutgoingWebhook, error)
	CreateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookCreateForm) (*model.OutgoingWebhook, error)
	UpdateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookUpdateForm) (*model.OutgoingWebhook, error)
	DeleteOutgoingWebhookByUser(uid uint, id uint) error
	DeleteOutgoingWebhooksByUser(uid uint, ids []uint) (int64, error)
}
