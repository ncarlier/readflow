package db

import "github.com/ncarlier/reader/pkg/model"

// APIKeyRepository is the repository interface to manage API keys
type APIKeyRepository interface {
	GetAPIKeyByID(id uint) (*model.APIKey, error)
	GetAPIKeyByToken(token string) (*model.APIKey, error)
	GetAPIKeyByUserIDAndAlias(userID uint, alias string) (*model.APIKey, error)
	GetAPIKeysByUserID(userID uint) ([]model.APIKey, error)
	CreateOrUpdateAPIKey(apiKey model.APIKey) (*model.APIKey, error)
	DeleteAPIKey(apiKey model.APIKey) error
	DeleteAPIKeys(uid uint, ids []uint) (int64, error)
}
