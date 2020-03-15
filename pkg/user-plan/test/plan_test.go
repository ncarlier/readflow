package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/assert"
	userplan "github.com/ncarlier/readflow/pkg/user-plan"
)

func TestNewEmptyUserPlans(t *testing.T) {
	userPlans, err := userplan.NewUserPlans("")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 0, len(userPlans.Plans), "unexpected number of plan")
	assert.True(t, userPlans.GetPlan("test") == nil, "plan should not be found")
}

func TestNewUserPlans(t *testing.T) {
	userPlans, err := userplan.NewUserPlans("user-plans.yml")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 2, len(userPlans.Plans), "unexpected number of plan")
	plan := userPlans.GetPlan("test")
	assert.Equal(t, "starter", plan.Name, "unexpected plan name")
	assert.Equal(t, uint(200), plan.TotalArticles, "unexpected total articles value")
	plan = userPlans.GetPlan("premium")
	assert.Equal(t, "premium", plan.Name, "unexpected plan name")
	assert.Equal(t, uint(2000), plan.TotalArticles, "unexpected total articles value")
}
