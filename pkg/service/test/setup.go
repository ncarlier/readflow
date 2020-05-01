package dbtest

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	userplan "github.com/ncarlier/readflow/pkg/user-plan"
)

var testDB db.DB

const defaultUsername = "test"

const defaultDBConnString = "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"

var testUser *model.User

var testContext context.Context

func assertUserExists(t *testing.T, username string) *model.User {
	user, err := testDB.GetUserByUsername(username)
	assert.Nil(t, err)
	if user != nil {
		return user
	}

	user = &model.User{
		Username: username,
	}
	user, err = testDB.CreateOrUpdateUser(*user)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, user.ID)
	assert.Equal(t, username, user.Username)
	return user
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	var err error
	testDB, err = db.NewDB(getEnv("DB", defaultDBConnString))
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	testUser = assertUserExists(t, defaultUsername)
	testContext = context.Background()
	testContext = context.WithValue(testContext, constant.UserID, *testUser.ID)
	userPlans, _ := userplan.NewUserPlans("user-plans.yml")
	conf := config.Config{}

	service.Configure(conf, testDB, userPlans)
	if err != nil {
		t.Fatalf("unable to setup service registry: %v", err)
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		defer testDB.Close()
		testDB.DeleteUser(*testUser)
	}
}

func init() {
	logger.Configure("debug", true, "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("READFLOW_" + key); ok {
		return value
	}
	return fallback
}
