package oidc

import (
	"crypto/rsa"
	"fmt"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

// Keystore OIDC keystore
type Keystore struct {
	conf  Configuration
	store sync.Map
}

// NewOIDCKeystore create a new OIDC keystore
func NewOIDCKeystore(conf Configuration) *Keystore {
	return &Keystore{
		conf: conf,
	}
}

// GetKey retrieve a key from the keystore
func (k *Keystore) GetKey(id string) (*rsa.PublicKey, error) {
	key, ok := k.store.Load(id)
	if ok {
		return key.(*rsa.PublicKey), nil
	}

	pubKey, err := k.getKey(id)
	if err == nil {
		k.store.Store(id, pubKey)
	}
	return pubKey, err
}

func (k *Keystore) getKey(id string) (*rsa.PublicKey, error) {
	var pem = ""
	for _, jwk := range k.conf.JSONWebKeySet {
		if jwk.Kid == id {
			var err error
			pem, err = jwkToPEM(jwk)
			if err != nil {
				return nil, err
			}
		}
	}

	if pem == "" {
		err := fmt.Errorf("unable to fing key #%s in OIDC configuration", id)
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
}
