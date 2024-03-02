package types

import (
	"encoding/hex"
)

type HexString struct {
	Value []byte
}

func (hs *HexString) UnmarshalText(text []byte) error {
	var err error
	hs.Value, err = hex.DecodeString(string(text))
	return err
}
