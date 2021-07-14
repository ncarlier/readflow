package helper

import (
	"github.com/speps/go-hashids/v2"
)

type HashIDHandler struct {
	provider *hashids.HashID
}

// NewHashIDHandler creates hashid handler
func NewHashIDHandler(salt string) (*HashIDHandler, error) {
	hd := hashids.NewData()
	hd.Salt = salt
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
