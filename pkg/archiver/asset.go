package archiver

import (
	"bytes"
	"encoding/gob"
)

// WebAsset is asset that used in a web page.
type WebAsset struct {
	Data        []byte
	ContentType string
}

// Encode web asset structure to byte array
func (obj WebAsset) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeWebAsset byte array to web asset sctructure
func DecodeWebAsset(b []byte) (*WebAsset, error) {
	obj := WebAsset{}
	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
