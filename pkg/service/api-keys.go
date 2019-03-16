package service

import (
	"github.com/ncarlier/reader/pkg/model"
)

// GetAPIKeyByToken returns an API key by its token
func (reg *Registry) GetAPIKeyByToken(token string) (*model.APIKey, error) {
	return reg.db.GetAPIKeyByToken(token)
}
