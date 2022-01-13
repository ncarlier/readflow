package dbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
)

func assertCategoryExists(t *testing.T, uid uint, title string, notif string) *model.Category {
	category, err := testDB.GetCategoryByUserAndTitle(uid, title)
	assert.Nil(t, err)
	if category != nil {
		return category
	}

	createForm := model.CategoryCreateForm{
		Title:                title,
		NotificationStrategy: notif,
	}

	category, err = testDB.CreateCategoryForUser(uid, createForm)
	assert.Nil(t, err)
	assert.NotNil(t, category)
	assert.NotNil(t, category.ID)
	assert.Equal(t, title, category.Title)
	assert.Equal(t, notif, category.NotificationStrategy)
	assert.Nil(t, category.Rule)
	return category
}
func TestCreateAndUpdateCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	uid := *testUser.ID
	title := "My test category"
	category := assertCategoryExists(t, uid, title, "none")

	// Update category title
	title = "My updated category"
	notif := "global"
	update := model.CategoryUpdateForm{
		ID:                   *category.ID,
		Title:                &title,
		NotificationStrategy: &notif,
	}
	category, err := testDB.UpdateCategoryForUser(uid, update)
	assert.Nil(t, err)
	assert.NotNil(t, category)
	assert.NotNil(t, category.ID)
	assert.Equal(t, title, category.Title)
	assert.Nil(t, category.Rule)
	assert.Equal(t, notif, category.NotificationStrategy)

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
	assert.Equal(t, notif, category.NotificationStrategy)

	// Count categories of test user
	nb, err := testDB.CountCategoriesByUser(uid)
	assert.Nil(t, err)
	assert.Positive(t, nb)
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	title := "My updated category"
	category := assertCategoryExists(t, uid, title, "none")

	categories, err := testDB.GetCategoriesByUser(uid)
	assert.Nil(t, err)
	assert.Positive(t, len(categories), "categories should not be empty")

	err = testDB.DeleteCategoryByUser(uid, *category.ID)
	assert.Nil(t, err)

	category, err = testDB.GetCategoryByUserAndTitle(uid, title)
	assert.Nil(t, err)
	assert.Nil(t, category)
}
