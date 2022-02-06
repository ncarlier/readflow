package sanitizer

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

// countLines count file lines
func countLines(r io.Reader) (uint, error) {
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

func open(location string) (io.ReadCloser, error) {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		resp, err := http.DefaultClient.Get(location)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	return os.Open(location)
}
