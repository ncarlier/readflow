package tooling

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash creats a hash from a payload string
func Hash(payload string) string {
	hasher := md5.New()
	hasher.Write([]byte(payload))
	return hex.EncodeToString(hasher.Sum(nil))
}
