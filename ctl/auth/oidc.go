package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// OIDCTokenCache stores OAuth tokens on disk
type OIDCTokenCache struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// oidcProviderConfig holds OIDC discovery endpoints
type oidcProviderConfig struct {
	TokenEndpoint               string `json:"token_endpoint"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
}

// deviceCodeResponse is the response from the device authorization endpoint
type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// tokenResponse is the response from the token endpoint
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error,omitempty"`
}

// OIDCTransport is an HTTP transport that adds OIDC Bearer token headers
type OIDCTransport struct {
	Issuer    string
	ClientID  string
	Transport http.RoundTripper

	cache    *OIDCTokenCache
	provider *oidcProviderConfig
}

// RoundTrip implements the http.RoundTripper interface
func (t *OIDCTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.ensureToken(); err != nil {
		return nil, fmt.Errorf("unable to obtain OIDC token: %w", err)
	}
	req.Header.Set("Authorization", t.cache.TokenType+" "+t.cache.AccessToken)
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return transport.RoundTrip(req)
}

func (t *OIDCTransport) ensureToken() error {
	// try loading cached token
	if t.cache == nil {
		cache, err := loadTokenCache()
		if err == nil && cache != nil {
			t.cache = cache
		}
	}

	// if we have a valid (non-expired) token, use it
	if t.cache != nil && time.Now().Before(t.cache.ExpiresAt) {
		return nil
	}

	// discover provider endpoints
	if t.provider == nil {
		provider, err := discoverOIDC(t.Issuer)
		if err != nil {
			return err
		}
		t.provider = provider
	}

	// try refreshing with refresh token
	if t.cache != nil && t.cache.RefreshToken != "" {
		if err := t.refreshToken(); err == nil {
			return nil
		}
	}

	// fall back to device code flow
	return t.deviceCodeFlow()
}

func discoverOIDC(issuer string) (*oidcProviderConfig, error) {
	wellKnown := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch OIDC discovery document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OIDC discovery returned HTTP %d", resp.StatusCode)
	}

	var config oidcProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode OIDC discovery document: %w", err)
	}
	return &config, nil
}

func (t *OIDCTransport) deviceCodeFlow() error {
	// Step 1: Request device code
	data := url.Values{
		"client_id": {t.ClientID},
		"scope":     {"openid"},
	}
	resp, err := http.PostForm(t.provider.DeviceAuthorizationEndpoint, data)
	if err != nil {
		return fmt.Errorf("unable to request device code: %w", err)
	}
	defer resp.Body.Close()

	var dcResp deviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&dcResp); err != nil {
		return fmt.Errorf("unable to decode device code response: %w", err)
	}

	// Step 2: Display verification URL and user code
	fmt.Fprintf(os.Stderr, "\nTo authenticate, visit: %s\n", dcResp.VerificationURI)
	fmt.Fprintf(os.Stderr, "And enter code: %s\n\n", dcResp.UserCode)
	if dcResp.VerificationURIComplete != "" {
		fmt.Fprintf(os.Stderr, "Or open: %s\n\n", dcResp.VerificationURIComplete)
	}

	// Step 3: Poll token endpoint
	interval := time.Duration(dcResp.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(dcResp.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(interval)

		tokenData := url.Values{
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
			"device_code": {dcResp.DeviceCode},
			"client_id":   {t.ClientID},
		}
		tokenResp, err := http.PostForm(t.provider.TokenEndpoint, tokenData)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(tokenResp.Body)
		tokenResp.Body.Close()
		if err != nil {
			continue
		}

		var tokResp tokenResponse
		if err := json.Unmarshal(body, &tokResp); err != nil {
			continue
		}

		switch tokResp.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
			continue
		case "":
			// success
			t.cache = &OIDCTokenCache{
				AccessToken:  tokResp.AccessToken,
				RefreshToken: tokResp.RefreshToken,
				TokenType:    tokResp.TokenType,
				ExpiresAt:    time.Now().Add(time.Duration(tokResp.ExpiresIn) * time.Second),
			}
			return saveTokenCache(t.cache)
		default:
			return fmt.Errorf("OIDC token error: %s", tokResp.Error)
		}
	}
	return fmt.Errorf("device code flow timed out")
}

func (t *OIDCTransport) refreshToken() error {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {t.cache.RefreshToken},
		"client_id":     {t.ClientID},
	}
	resp, err := http.PostForm(t.provider.TokenEndpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokResp); err != nil {
		return err
	}
	if tokResp.Error != "" {
		return fmt.Errorf("refresh error: %s", tokResp.Error)
	}

	t.cache = &OIDCTokenCache{
		AccessToken:  tokResp.AccessToken,
		RefreshToken: tokResp.RefreshToken,
		TokenType:    tokResp.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(tokResp.ExpiresIn) * time.Second),
	}
	if t.cache.RefreshToken == "" {
		// keep old refresh token if new one not provided
		t.cache.RefreshToken = data.Get("refresh_token")
	}
	return saveTokenCache(t.cache)
}

func tokenCachePath() string {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			cacheDir = "/tmp"
		} else {
			cacheDir = filepath.Join(home, ".cache")
		}
	}
	return filepath.Join(cacheDir, "readflow-ctl", "tokens.json")
}

func loadTokenCache() (*OIDCTokenCache, error) {
	path := tokenCachePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cache OIDCTokenCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func saveTokenCache(cache *OIDCTokenCache) error {
	path := tokenCachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
