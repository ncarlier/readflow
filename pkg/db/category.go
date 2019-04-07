package db

import "github.com/ncarlier/readflow/pkg/model"

// CategoryRepository is the repository interface to manage categories
type CategoryRepository interface {
	GetCategoryByID(id uint) (*model.Category, error)
	GetCategoryByUserIDAndTitle(uid uint, title string) (*model.Category, error)
	GetCategoriesByUserID(uid uint) ([]model.Category, error)
	CreateOrUpdateCategory(category model.Category) (*model.Category, error)
	DeleteCategory(category model.Category) error
	DeleteCategories(uid uint, ids []uint) (int64, error)
}
