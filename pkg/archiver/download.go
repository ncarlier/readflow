package archiver

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/constant"
)

func (arc *WebArchiver) downloadFile(url string, parentURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", constant.UserAgent)
	if parentURL != "" {
		req.Header.Set("Referer", parentURL)
	}

	resp, err := arc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
