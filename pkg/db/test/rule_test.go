package dbtest

import (
	"testing"

	"github.com/ncarlier/reader/pkg/assert"
	"github.com/ncarlier/reader/pkg/model"
)

func assertRuleExists(t *testing.T, rule model.Rule) *model.Rule {
	result, err := testDB.GetRuleByUserIDAndAlias(*rule.UserID, rule.Alias)
	assert.Nil(t, err, "error on getting rule by user and alias should be nil")
	if result != nil {
		id := *result.ID
		rule.ID = &id
	}

	result, err = testDB.CreateOrUpdateRule(rule)
	assert.Nil(t, err, "error on create/update rule should be nil")
	assert.NotNil(t, result, "rule shouldn't be nil")
	assert.NotNil(t, result.ID, "rule ID shouldn't be nil")
	assert.Equal(t, *rule.UserID, *result.UserID, "")
	assert.Equal(t, *rule.CategoryID, *result.CategoryID, "")
	assert.Equal(t, rule.Alias, result.Alias, "")
	assert.Equal(t, rule.Rule, result.Rule, "")
	assert.Equal(t, rule.Priority, result.Priority, "")
	return result
}
func TestCreateOrUpdateRule(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	category := assertCategoryExists(t, testUser.ID, "My category")

	rule := model.Rule{
		Alias:      "My rule",
		UserID:     testUser.ID,
		CategoryID: category.ID,
		Rule:       "title = \"foo\"",
		Priority:   0,
	}

	// Create rule
	newRule := assertRuleExists(t, rule)

	newRule.Alias = "My updated rule"

	// Update rule
	updatedRule, err := testDB.CreateOrUpdateRule(*newRule)
	assert.Nil(t, err, "error on update rule should be nil")
	assert.NotNil(t, updatedRule, "rule shouldn't be nil")
	assert.NotNil(t, updatedRule.ID, "rule ID shouldn't be nil")
	assert.Equal(t, newRule.Alias, updatedRule.Alias, "")

	// Cleanup
	err = testDB.DeleteRule(*updatedRule)
	assert.Nil(t, err, "error on cleanup should be nil")
}

func TestDeleteRule(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	// Create category
	category := assertCategoryExists(t, testUser.ID, "My category")

	alias := "My rule"

	// Assert rule exists
	rule := &model.Rule{
		Alias:      alias,
		UserID:     testUser.ID,
		CategoryID: category.ID,
		Rule:       "title = \"foo\"",
		Priority:   0,
	}
	rule = assertRuleExists(t, *rule)

	rules, err := testDB.GetRulesByUserID(*testUser.ID)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(rules) > 0, "rules should not be empty")

	err = testDB.DeleteRule(*rule)
	assert.Nil(t, err, "error on delete should be nil")

	rule, err = testDB.GetRuleByUserIDAndAlias(*testUser.ID, alias)
	assert.Nil(t, err, "error should be nil")
	assert.True(t, rule == nil, "rule should be nil")
}
