package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/secret"
	"github.com/stretchr/testify/require"
)

func TestSealAndUnSealSecrets(t *testing.T) {
	engine, err := secret.NewSecretsEngineProvider("file://./secret.key")
	require.Nil(t, err)
	require.NotNil(t, engine)

	secrets := make(secret.Secrets)
	secrets["foo"] = "bar"
	secrets["zoo"] = ""

	err = engine.Apply(secret.Seal, &secrets)
	require.Nil(t, err)
	require.NotEqual(t, "bar", secrets["foo"])
	require.Empty(t, secrets["zoo"])

	err = engine.Apply(secret.UnSeal, &secrets)
	require.Nil(t, err)
	require.Equal(t, "bar", secrets["foo"])
	require.Empty(t, secrets["zoo"])
}
