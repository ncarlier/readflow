package oidc

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/ncarlier/readflow/pkg/defaults"
)

// Keystore OIDC keystore
type Keystore struct {
	conf  *Configuration
	store sync.Map
	lock  sync.Mutex
}

// NewOIDCKeystore create a new OIDC keystore
func NewOIDCKeystore(conf *Configuration) (*Keystore, error) {
	ks := &Keystore{
		conf: conf,
	}
	if err := ks.fetch(); err != nil {
		return nil, err
	}
	return ks, nil
}

func (k *Keystore) fetch() error {
	k.lock.Lock()
	defer k.lock.Unlock()

	resp, err := defaults.HTTPClient.Get(k.conf.JwksURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var jwks = JSONWebKeySet{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return err
	}

	for _, jwk := range jwks.Keys {
		if pem, err := jwkToPEM(jwk); err == nil {
			if pubkey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem)); err == nil {
				k.store.Store(jwk.Kid, pubkey)
			}
		}
	}

	return nil
}

// GetKey retrieve a key from the keystore
func (k *Keystore) GetKey(id string) (*rsa.PublicKey, error) {
	if key, ok := k.store.Load(id); ok {
		return key.(*rsa.PublicKey), nil
	}

	// refresh store in case of key rotation...
	if err := k.fetch(); err != nil {
		return nil, err
	}
	if key, ok := k.store.Load(id); ok {
		return key.(*rsa.PublicKey), nil
	}
	return nil, fmt.Errorf("unable to fing key #%s in OIDC configuration", id)
}
