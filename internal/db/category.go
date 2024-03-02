package db

import "github.com/ncarlier/readflow/internal/model"

// CategoryRepository is the repository interface to manage categories
type CategoryRepository interface {
	GetCategoryByID(id uint) (*model.Category, error)
	GetCategoryByUserAndTitle(uid uint, title string) (*model.Category, error)
	GetCategoriesByUser(uid uint) ([]model.Category, error)
	CountCategoriesByUser(uid uint) (uint, error)
	CreateCategoryForUser(uid uint, form model.CategoryCreateForm) (*model.Category, error)
	UpdateCategoryForUser(uid uint, form model.CategoryUpdateForm) (*model.Category, error)
	DeleteCategoryByUser(uid uint, ID uint) error
	DeleteCategoriesByUser(uid uint, ids []uint) (int64, error)
}
