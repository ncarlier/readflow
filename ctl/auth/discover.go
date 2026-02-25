package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// serverInfo mirrors the /info response from the readflow server.
type serverInfo struct {
	Authority string `json:"authority"`
}

// DiscoverAuthMethod calls GET <endpoint>/info and deduces the auth method:
//   - authority == ""        → method "basic"
//   - authority is a URL    → method "oidc", issuer = authority
func DiscoverAuthMethod(endpoint string) (method, issuer string, err error) {
	infoURL := strings.TrimRight(endpoint, "/") + "/info"
	resp, err := http.Get(infoURL)
	if err != nil {
		return "", "", fmt.Errorf("unable to reach %s: %w", infoURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GET %s returned HTTP %d", infoURL, resp.StatusCode)
	}

	var info serverInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", "", fmt.Errorf("unable to decode /info response: %w", err)
	}

	if info.Authority == "" {
		return "basic", "", nil
	}

	if _, err := url.ParseRequestURI(info.Authority); err != nil {
		return "", "", fmt.Errorf("authority field %q is not a valid URL: %w", info.Authority, err)
	}

	return "oidc", info.Authority, nil
}
