package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertCategoryExists(t *testing.T, uid uint, title string) *model.Category {
	category, err := testDB.GetCategoryByUserAndTitle(uid, title)
	assert.Nil(t, err)
	if category != nil {
		return category
	}

	createForm := model.CategoryCreateForm{
		Title: title,
	}

	category, err = testDB.CreateCategoryForUser(uid, createForm)
	assert.Nil(t, err)
	assert.NotNil(t, category)
	assert.NotNil(t, category.ID)
	assert.Equal(t, title, category.Title)
	assert.Nil(t, category.Rule)
	return category
}
func TestCreateAndUpdateCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	uid := *testUser.ID
	title := "My test category"
	category := assertCategoryExists(t, uid, title)

	// Update category title
	title = "My updated category"
	update := model.CategoryUpdateForm{
		ID:    *category.ID,
		Title: &title,
	}
	category, err := testDB.UpdateCategoryForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, category)
	assert.NotNil(t, category.ID)
	assert.Equal(t, title, category.Title)
	assert.Nil(t, category.Rule)

	// Update category title
	rule := "title matches \"test\""
	update = model.CategoryUpdateForm{
		ID:   *category.ID,
		Rule: &rule,
	}
	category, err = testDB.UpdateCategoryForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, category)
	assert.NotNil(t, category.ID)
	assert.Equal(t, title, category.Title)
	assert.Equal(t, rule, *category.Rule)

	// Count categories of test user
	nb, err := testDB.CountCategoriesByUser(uid)
	assert.Nil(t, err)
	assert.True(t, nb >= 0)
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	title := "My updated category"
	category := assertCategoryExists(t, uid, title)

	categories, err := testDB.GetCategoriesByUser(uid)
	assert.Nil(t, err)
	assert.True(t, len(categories) > 0, "categories should not be empty")

	err = testDB.DeleteCategoryByUser(uid, *category.ID)
	assert.Nil(t, err)

	category, err = testDB.GetCategoryByUserAndTitle(uid, title)
	assert.Nil(t, err)
	assert.Nil(t, category)
}
