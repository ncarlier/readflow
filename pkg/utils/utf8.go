package utils

import (
	"bytes"
	"io"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
)

// ToUTF8 converts ISO string to UTF8
func ToUTF8(iso string) string {
	if utf8.ValidString(iso) {
		return iso
	}
	isoBuf := []byte(iso)
	buf := make([]rune, len(isoBuf))
	for i, b := range isoBuf {
		buf[i] = rune(b)
	}
	return string(buf)
}

// NewUTF8Reader converts a reader from a charset to UTF-8
func NewUTF8Reader(reader io.Reader, sourceCharset string) (io.Reader, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(b)
	return charset.NewReaderLabel(sourceCharset, br)
}
