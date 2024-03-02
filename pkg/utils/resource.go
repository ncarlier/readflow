package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// OpenResource open local or remote resource as a Reader
func OpenResource(location string) (io.ReadCloser, error) {
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
