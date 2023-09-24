package oidc

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Keystore OIDC keystore
type Keystore struct {
	conf   Configuration
	keys   []JSONWebKey
	store  sync.Map
	ticker *time.Ticker
}

// NewOIDCKeystore create a new OIDC keystore
func NewOIDCKeystore(conf Configuration) (*Keystore, error) {
	ks := &Keystore{
		conf: conf,
	}
	if err := ks.refresh(); err != nil {
		return nil, err
	}
	return ks, nil
}

func (k *Keystore) refresh() error {
	resp, err := http.Get(k.conf.JwksURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&k.keys)
	if err != nil {
		return err
	}
	return nil
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
	for _, jwk := range k.keys {
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

// Start keystore periodic refresh job
func (k *Keystore) Start() {
	if k.ticker != nil {
		return
	}
	k.ticker = time.NewTicker(time.Hour)
	for range k.ticker.C {
		k.refresh()
	}
}

// Stop keystore periodic refresh job
func (k *Keystore) Stop() {
	k.ticker.Stop()
}
