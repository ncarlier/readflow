package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/ncarlier/readflow/pkg/helper"
)

type fileSecretProvider struct {
	key []byte
}

func newLocalSecretsEngineProvider(uri string) (EngineProvider, error) {
	input, err := helper.OpenResource(uri)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	buf := make([]byte, 32)
	_, err = input.Read(buf)
	if err != nil {
		return nil, err
	}

	return &fileSecretProvider{
		key: buf,
	}, nil
}

func (p fileSecretProvider) Seal(secrets *Secrets) error {
	blockCipher, err := aes.NewCipher(p.key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return err
	}

	for k, v := range *secrets {
		if v == "" {
			continue
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		(*secrets)[k] = string(ciphertext)
	}

	return nil
}

func (p fileSecretProvider) UnSeal(secrets *Secrets) error {
	blockCipher, err := aes.NewCipher(p.key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return err
	}

	for k, v := range *secrets {
		if v == "" {
			continue
		}
		data := []byte(v)
		nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}
		(*secrets)[k] = string(plaintext)
	}

	return nil
}
