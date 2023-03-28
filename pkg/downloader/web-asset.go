package downloader

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

// WebAsset is a structure used to store file properties.
type WebAsset struct {
	Data        []byte
	ContentType string
	Name        string
}

// Encode file asset structure to byte array
func (obj WebAsset) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToDataURL returns base64 encoded data URL of the file asset
func (obj WebAsset) ToDataURL() string {
	b64encoded := base64.StdEncoding.EncodeToString(obj.Data)
	return fmt.Sprintf("data:%s;base64,%s", obj.ContentType, b64encoded)
}

// NewWebAsset byte array to file asset sctructure
func NewWebAsset(b []byte) (*WebAsset, error) {
	obj := WebAsset{}
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
