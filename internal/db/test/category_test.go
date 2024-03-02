package dbtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
)

func assertCategoryExists(t *testing.T, uid uint, title, notif string) *model.Category {
	category, err := testDB.GetCategoryByUserAndTitle(uid, title)
	require.Nil(t, err)
	if category != nil {
		return category
	}

	createForm := model.CategoryCreateForm{
		Title: title,
	}

	category, err = testDB.CreateCategoryForUser(uid, createForm)
	require.Nil(t, err)
	require.NotNil(t, category)
	require.NotNil(t, category.ID)
	require.Equal(t, title, category.Title)
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
	update := model.CategoryUpdateForm{
		ID:    *category.ID,
		Title: &title,
	}
	category, err := testDB.UpdateCategoryForUser(uid, update)
	require.Nil(t, err)
	require.NotNil(t, category)
	require.NotNil(t, category.ID)
	require.Equal(t, title, category.Title)

	// Count categories of test user
	nb, err := testDB.CountCategoriesByUser(uid)
	require.Nil(t, err)
	require.Positive(t, nb)
}

func TestDeleteCategory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Assert category exists
	uid := *testUser.ID
	title := "My updated category"
	category := assertCategoryExists(t, uid, title, "none")

	categories, err := testDB.GetCategoriesByUser(uid)
	require.Nil(t, err)
	require.Positive(t, len(categories), "categories should not be empty")

	err = testDB.DeleteCategoryByUser(uid, *category.ID)
	require.Nil(t, err)

	category, err = testDB.GetCategoryByUserAndTitle(uid, title)
	require.Nil(t, err)
	require.Nil(t, category)
}
