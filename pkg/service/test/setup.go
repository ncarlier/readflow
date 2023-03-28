package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
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
	conf := config.NewConfig()
	if err := conf.LoadFile("test.toml"); err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	if conf.Global.DatabaseURI == "" {
		conf.Global.DatabaseURI = defaultDBConnString
	}

	var err error
	testDB, err = db.NewDB(conf.Global.DatabaseURI)
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}
	testUser = assertUserExists(t, defaultUsername)
	testContext = context.Background()
	testContext = context.WithValue(testContext, constant.ContextUserID, *testUser.ID)
	downloadCache, _ := cache.NewDefault("")

	service.Configure(*conf, testDB, downloadCache)
	if err != nil {
		t.Fatalf("unable to setup service registry: %v", err)
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		defer testDB.Close()
		defer downloadCache.Close()
		testDB.DeleteUser(*testUser)
	}
}

func init() {
	logger.Configure("debug", true, "")
}
