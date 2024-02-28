package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/oidc"
	"github.com/stretchr/testify/require"
)

const issuer = "https://accounts.readflow.app"

func TestOIDCDiscovery(t *testing.T) {
	client, err := oidc.NewOIDCClient(issuer, "", "")
	require.Nil(t, err)
	require.Equal(t, issuer, client.Config.Issuer)
	require.NotNil(t, client.Keystore)
}
