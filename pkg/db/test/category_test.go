package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func assertCategoryExists(t *testing.T, userID *uint32, title string) *model.Category {
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

	// Assert user exists
	user := assertUserExists(t, "test-002")

	assertCategoryExists(t, user.ID, title)
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	title := "My test category"

	// Assert user exists
	user := assertUserExists(t, "test-002")

	// Assert category exists
	category := assertCategoryExists(t, user.ID, title)

	err := testDB.DeleteCategory(*category)
	assert.Nil(t, err, "error on delete should be nil")

	categories, err := testDB.GetCategoriesByUserID(*user.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(categories) == 0, "categories should be empty")
}
