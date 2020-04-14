package db

import "github.com/ncarlier/readflow/pkg/model"

// APIKeyRepository is the repository interface to manage API keys
type APIKeyRepository interface {
	GetAPIKeyByID(id uint) (*model.APIKey, error)
	GetAPIKeyByToken(token string) (*model.APIKey, error)
	GetAPIKeyByUserAndAlias(uid uint, alias string) (*model.APIKey, error)
	GetAPIKeysByUser(uid uint) ([]model.APIKey, error)
	CreateAPIKeyForUser(uid uint, form model.APIKeyCreateForm) (*model.APIKey, error)
	UpdateAPIKeyForUser(uid uint, form model.APIKeyUpdateForm) (*model.APIKey, error)
	DeleteAPIKeyByUser(uid uint, id uint) error
	DeleteAPIKeysByUser(uid uint, ids []uint) (int64, error)
}
