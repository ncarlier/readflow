package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	conf := config.NewConfig()
	assert.Equal(t, ":8080", conf.HTTP.ListenAddr)
	assert.Empty(t, conf.ImageProxy.URL)
	assert.Nil(t, conf.GetUserPlan("test"), "plan should not be found")
}

func TestLaodConfigFromFile(t *testing.T) {
	conf := config.NewConfig()
	err := conf.LoadFile("test.toml")
	require.Nil(t, err)
	// Default overide
	assert.Equal(t, "localhost:8081", conf.HTTP.ListenAddr)
	// Env variable substitution
	assert.Equal(t, []byte("test"), conf.Hash.SecretSalt.Value)
	// Default if empty
	assert.Equal(t, "http://localhost:8080", conf.HTTP.PublicURL)
	// Sub attribute
	assert.Equal(t, "https://1..9:1..9@sentry.io/1..9", conf.Integration.Sentry.DSN)
}

func TestUserPlans(t *testing.T) {
	conf := config.NewConfig()
	err := conf.LoadFile("test.toml")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(conf.UserPlans), "unexpected number of plan")
	plan := conf.GetUserPlan("test")
	assert.Equal(t, "starter", plan.Name)
	assert.Equal(t, uint(200), plan.ArticlesLimit, "unexpected articles limit value")
	plan = conf.GetUserPlan("premium")
	assert.Equal(t, "premium", plan.Name)
	assert.Equal(t, uint(2000), plan.ArticlesLimit, "unexpected articles limit value")
}
