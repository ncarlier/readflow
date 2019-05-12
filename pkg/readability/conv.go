package readability

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/net/html/charset"
)

// NewUTF8Reader converts a reader from a charset to UTF-8
func NewUTF8Reader(reader io.Reader, sourceCharset string) (io.Reader, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(b)
	return charset.NewReaderLabel(sourceCharset, br)
}
