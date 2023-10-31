package config

import (
	"encoding/hex"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type hex_string struct {
	Value []byte
}

func (hs *hex_string) UnmarshalText(text []byte) error {
	var err error
	hs.Value, err = hex.DecodeString(string(text))
	return err
}
