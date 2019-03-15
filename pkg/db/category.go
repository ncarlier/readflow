package db

import "github.com/ncarlier/reader/pkg/model"

// CategoryRepository is the repository interface to manage categories
type CategoryRepository interface {
	GetCategoryByUserIDAndTitle(userID uint32, title string) (*model.Category, error)
	GetCategoriesByUserID(userID uint32) ([]model.Category, error)
	CreateOrUpdateCategory(category model.Category) (*model.Category, error)
	DeleteCategory(category model.Category) error
}
