package model

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

// FileAsset is a structure used to store file properties.
type FileAsset struct {
	Data        []byte
	ContentType string
	Name        string
}

// Encode file asset structure to byte array
func (obj FileAsset) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToDataURL returns base64 encoded data URL of the file asset
func (obj FileAsset) ToDataURL() string {
	b64encoded := base64.StdEncoding.EncodeToString(obj.Data)
	return fmt.Sprintf("data:%s;base64,%s", obj.ContentType, b64encoded)
}

// DecodeFileAsset byte array to file asset sctructure
func DecodeFileAsset(b []byte) (*FileAsset, error) {
	obj := FileAsset{}
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
