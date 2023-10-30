package helper

import (
	"crypto/hmac"
	"encoding/base64"
	"hash"
)

// Sign value with secret
func Sign(algo func() hash.Hash, value, secret string, truncate int) string {
	h := hmac.New(algo, []byte(secret))
	h.Write([]byte(value))
	sig := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if truncate > 0 && len(sig) > truncate {
		return sig[:truncate]
	}
	return sig
}
