package dbtest

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/assert"
	"github.com/ncarlier/readflow/pkg/model"
)

func assertCategoryExists(t *testing.T, userID *uint, title string) *model.Category {
	category, err := testDB.GetCategoryByUserIDAndTitle(*userID, title)
	assert.Nil(t, err, "error on getting category by user and title should be nil")
	if category != nil {
		return category
	}

	category = &model.Category{
		UserID: userID,
		Title:  title,
	}

	category, err = testDB.CreateOrUpdateCategory(*category)
	assert.Nil(t, err, "error on create/update category should be nil")
	assert.NotNil(t, category, "category shouldn't be nil")
	assert.NotNil(t, category.ID, "category ID shouldn't be nil")
	assert.Equal(t, title, category.Title, "")
	return category
}
func TestCreateOrUpdateCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	title := "My test category"

	// Create category
	category := assertCategoryExists(t, testUser.ID, title)

	title = "My updated category"
	category.Title = title

	// Update category
	category, err := testDB.CreateOrUpdateCategory(*category)
	assert.Nil(t, err, "error on update category should be nil")
	assert.NotNil(t, category, "category shouldn't be nil")
	assert.NotNil(t, category.ID, "category ID shouldn't be nil")
	assert.Equal(t, title, category.Title, "")

	nb, err := testDB.CountCategoriesByUserID(*testUser.ID)
	assert.Nil(t, err, "")
	assert.True(t, nb >= 0, "")
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	title := "My updated category"

	// Assert category exists
	category := assertCategoryExists(t, testUser.ID, title)

	categories, err := testDB.GetCategoriesByUserID(*testUser.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(categories) > 0, "categories should not be empty")

	err = testDB.DeleteCategory(*category)
	assert.Nil(t, err, "error on delete should be nil")

	category, err = testDB.GetCategoryByUserIDAndTitle(*testUser.ID, title)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, category == nil, "category should be nil")
}
