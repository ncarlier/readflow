package utils

import (
	"bytes"
	"io"
)

// CountLines count lines of a Reader
func CountLines(r io.Reader) (uint, error) {
	buf := make([]byte, 32*1024)
	count := uint(0)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += uint(bytes.Count(buf[:c], lineSep))

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
