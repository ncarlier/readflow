package downloader

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ncarlier/readflow/pkg/utils"
)

// WebAsset is a structure used to store file properties.
type WebAsset struct {
	Data        []byte
	ContentType string
	Name        string
}

// Encode file asset structure to byte array
func (wa *WebAsset) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wa)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToDataURL returns base64 encoded data URL of the file asset
func (obj *WebAsset) ToDataURL() string {
	b64encoded := base64.StdEncoding.EncodeToString(obj.Data)
	return fmt.Sprintf("data:%s;base64,%s", obj.ContentType, b64encoded)
}

// Write to HTTP writer
func (wa *WebAsset) Write(w http.ResponseWriter, header http.Header) (int, error) {
	for k, vv := range header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("Content-Type", wa.ContentType)
	length := strconv.Itoa(len(wa.Data))
	if values := header.Values("Transfer-Encoding"); utils.ContainsString(values, "chunked") {
		// HACK: no Content-Length because of Transfer-Encoding=chunked
		w.Header().Set("X-Content-Length", length)
	} else {
		w.Header().Set("Content-Length", length)
	}
	w.Header().Set("Content-Disposition", "inline; filename=\""+wa.Name+"\"")
	return w.Write(wa.Data)
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
