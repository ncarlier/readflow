package tooling

import "unicode/utf8"

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
