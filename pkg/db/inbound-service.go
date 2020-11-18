package db

import "github.com/ncarlier/readflow/pkg/model"

// InboundServiceRepository is the repository interface to manage API keys
type InboundServiceRepository interface {
	GetInboundServiceByID(id uint) (*model.InboundService, error)
	GetInboundServiceByToken(token string) (*model.InboundService, error)
	GetInboundServiceByUserAndAlias(uid uint, alias string) (*model.InboundService, error)
	GetInboundServicesByUser(uid uint) ([]model.InboundService, error)
	CreateInboundServiceForUser(uid uint, form model.InboundServiceCreateForm) (*model.InboundService, error)
	UpdateInboundServiceForUser(uid uint, form model.InboundServiceUpdateForm) (*model.InboundService, error)
	DeleteInboundServiceByUser(uid uint, id uint) error
	DeleteInboundServicesByUser(uid uint, ids []uint) (int64, error)
}
