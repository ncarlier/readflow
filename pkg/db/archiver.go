package db

import "github.com/ncarlier/readflow/pkg/model"

// ArchiverRepository is the repository interface to manage Archivers
type ArchiverRepository interface {
	GetArchiverByID(id uint) (*model.Archiver, error)
	GetArchiverByUserAndAlias(uid uint, alias *string) (*model.Archiver, error)
	GetArchiversByUser(uid uint) ([]model.Archiver, error)
	CreateArchiverForUser(uid uint, form model.ArchiverCreateForm) (*model.Archiver, error)
	UpdateArchiverForUser(uid uint, form model.ArchiverUpdateForm) (*model.Archiver, error)
	DeleteArchiverByUser(uid uint, id uint) error
	DeleteArchiversByUser(uid uint, ids []uint) (int64, error)
}
