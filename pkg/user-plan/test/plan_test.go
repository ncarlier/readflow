package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	userplan "github.com/ncarlier/readflow/pkg/user-plan"
)

func TestNewEmptyUserPlans(t *testing.T) {
	userPlans, err := userplan.NewUserPlans("")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(userPlans.Plans))
	assert.Nil(t, userPlans.GetPlan("test"), "plan should not be found")
}

func TestNewUserPlans(t *testing.T) {
	userPlans, err := userplan.NewUserPlans("user-plans.yml")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(userPlans.Plans), "unexpected number of plan")
	plan := userPlans.GetPlan("test")
	assert.Equal(t, "starter", plan.Name)
	assert.Equal(t, uint(200), plan.TotalArticles, "unexpected total articles value")
	plan = userPlans.GetPlan("premium")
	assert.Equal(t, "premium", plan.Name)
	assert.Equal(t, uint(2000), plan.TotalArticles, "unexpected total articles value")
}
