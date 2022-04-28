package sanitizer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("invalid location: %s", location)
	}
	switch u.Scheme {
	case "http", "https":
		resp, err := http.DefaultClient.Get(location)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	case "file":
		return os.Open(u.Host + u.Path)
	}
	return nil, fmt.Errorf("unable to open file, scheme not supported: %s", location)
}
