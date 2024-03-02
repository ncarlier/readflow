package hashid

import (
	"github.com/speps/go-hashids/v2"
)

// HashIDHandler is used to hash a string with hashid algorythm
type HashIDHandler struct {
	provider *hashids.HashID
}

// NewHashIDHandler creates hashid handler
func NewHashIDHandler(salt []byte) (*HashIDHandler, error) {
	hd := hashids.NewData()
	hd.Salt = string(salt)
	provider, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return &HashIDHandler{
		provider: provider,
	}, nil
}

// Encode values into hashid
func (hid *HashIDHandler) Encode(values []int) string {
	result, _ := hid.provider.Encode(values)
	return result
}

// Decode value from hashid
func (hid *HashIDHandler) Decode(hash string) ([]int, error) {
	return hid.provider.DecodeWithError(hash)
}
