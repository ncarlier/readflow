package helper

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Hash creats a hash from a payload string
func Hash(values ...string) string {
	payload := strings.Join(values, "")
	hasher := md5.New()
	hasher.Write([]byte(payload))
	return hex.EncodeToString(hasher.Sum(nil))
}
