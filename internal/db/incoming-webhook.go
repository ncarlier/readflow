package db

import "github.com/ncarlier/readflow/internal/model"

// IncomingWebhookRepository is the repository interface to manage incoming webhooks
type IncomingWebhookRepository interface {
	GetIncomingWebhookByID(id uint) (*model.IncomingWebhook, error)
	GetIncomingWebhookByToken(token string) (*model.IncomingWebhook, error)
	GetIncomingWebhookByUserAndAlias(uid uint, alias string) (*model.IncomingWebhook, error)
	GetIncomingWebhooksByUser(uid uint) ([]model.IncomingWebhook, error)
	CountIncomingWebhooksByUser(uid uint) (uint, error)
	CreateIncomingWebhookForUser(uid uint, form model.IncomingWebhookCreateForm) (*model.IncomingWebhook, error)
	UpdateIncomingWebhookForUser(uid uint, form model.IncomingWebhookUpdateForm) (*model.IncomingWebhook, error)
	DeleteIncomingWebhookByUser(uid uint, id uint) error
	DeleteIncomingWebhooksByUser(uid uint, ids []uint) (int64, error)
}
