package db

import "github.com/ncarlier/reader/pkg/model"

// CategoryRepository is the repository interface to manage categories
type CategoryRepository interface {
	GetCategoryByID(id uint) (*model.Category, error)
	GetCategoryByUserIDAndTitle(userID uint, title string) (*model.Category, error)
	GetCategoriesByUserID(userID uint) ([]model.Category, error)
	CreateOrUpdateCategory(category model.Category) (*model.Category, error)
	DeleteCategory(category model.Category) error
}
