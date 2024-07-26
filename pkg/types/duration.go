package types

import (
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	val := string(text)
	if val == "" {
		val = "0s"
	}
	var err error
	d.Duration, err = time.ParseDuration(val)
	return err
}
