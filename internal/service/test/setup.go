package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/db"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/logger"
)

var testDB db.DB

const defaultUsername = "test"

const defaultDBConnString = "postgres://postgres:testpwd@localhost/readflow_test?sslmode=disable"

var testUser *model.User

var testContext context.Context

func requireUserExists(t *testing.T, username string) *model.User {
	user, err := testDB.GetUserByUsername(username)
	require.Nil(t, err)
	if user != nil {
		return user
	}

	user = &model.User{
		Username: username,
	}
	user, err = testDB.CreateOrUpdateUser(*user)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.ID)
	require.Equal(t, username, user.Username)
	return user
}

// SetupTestCase setup service test env
func SetupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	conf := config.NewConfig()
	if err := conf.LoadFile("test.toml"); err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	if conf.Database.URI == "" {
		conf.Database.URI = defaultDBConnString
	}

	var err error
	testDB, err = db.NewDB(conf.Database.URI)
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	testUser = requireUserExists(t, defaultUsername)
	testContext = context.Background()
	testContext = context.WithValue(testContext, global.ContextUserID, *testUser.ID)

	service.Configure(*conf, testDB)
	if err != nil {
		t.Fatalf("unable to setup service registry: %v", err)
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		defer service.Shutdown()
		defer testDB.Close()
		testDB.DeleteUser(*testUser)
	}
}

// GetTestContext returns test context
func GetTestContext() context.Context {
	return testContext
}

// GetTestUser returns test user
func GetTestUser() *model.User {
	return testUser
}

func init() {
	logger.Configure("debug", "text", "")
}
