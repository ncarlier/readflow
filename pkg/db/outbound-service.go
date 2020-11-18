package db

import "github.com/ncarlier/readflow/pkg/model"

// OutboundServiceRepository is the repository interface to manage OutboundServices
type OutboundServiceRepository interface {
	GetOutboundServiceByID(id uint) (*model.OutboundService, error)
	GetOutboundServiceByUserAndAlias(uid uint, alias *string) (*model.OutboundService, error)
	GetOutboundServicesByUser(uid uint) ([]model.OutboundService, error)
	CreateOutboundServiceForUser(uid uint, form model.OutboundServiceCreateForm) (*model.OutboundService, error)
	UpdateOutboundServiceForUser(uid uint, form model.OutboundServiceUpdateForm) (*model.OutboundService, error)
	DeleteOutboundServiceByUser(uid uint, id uint) error
	DeleteOutboundServicesByUser(uid uint, ids []uint) (int64, error)
}
