package test

import (
	"testing"

	"github.com/ncarlier/readflow/pkg/oidc"
	"github.com/stretchr/testify/require"
)

const issuer = "https://accounts.readflow.app"

func TestOIDCDiscovery(t *testing.T) {
	cfg, err := oidc.GetOIDCConfiguration(issuer)
	require.Nil(t, err)
	require.Equal(t, issuer, cfg.Issuer)

	keystore, err := oidc.NewOIDCKeystore(cfg)
	require.Nil(t, err)
	require.NotNil(t, keystore)

}