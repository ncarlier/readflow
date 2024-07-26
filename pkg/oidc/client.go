package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Cient of an OIDC issuer
type Client struct {
	client_id     string
	client_secret string
	Config        *Configuration
	Keystore      *Keystore
}

// NewOIDCClient create OpenID client from auto-discovery issuer endpoint
func NewOIDCClient(issuer, client_id, client_secret string) (*Client, error) {
	var config = &Configuration{}

	resp, err := http.Get(issuer + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid OIDC discovery response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(config)
	if err != nil {
		return nil, err
	}

	keystore, err := NewOIDCKeystore(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client_id:     client_id,
		client_secret: client_secret,
		Config:        config,
		Keystore:      keystore,
	}, nil
}

// Introspect call token introspection endpoint
func (c *Client) Introspect(token string) (*IntrospectionResponse, error) {
	form := url.Values{}
	form.Set("token", token)

	req, err := http.NewRequest("POST", c.Config.IntrospectionEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.client_id, c.client_secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to introspect token: %w", decodeErrorResponse(resp))
	}

	payload := &IntrospectionResponse{}
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// UserInfo call UserInfo endpoint
func (c *Client) UserInfo(token string) (*UserInfoResponse, error) {
	req, err := http.NewRequest("GET", c.Config.UserinfoEndpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to GET UserInfo: %w", decodeErrorResponse(resp))
	}

	payload := &UserInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// GetAuthorizationEndpoint return authorization endpoint
func (c *Client) GetAuthorizationEndpoint(redirectURI string) *url.URL {
	u, _ := url.Parse(c.Config.AuthorizationEndpoint)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("scope", "openid")
	//q.Set("state", "TODO")
	q.Set("client_id", c.client_id)
	if redirectURI != "" {
		q.Set("redirect_uri", redirectURI)
	}
	return u
}

func decodeErrorResponse(resp *http.Response) error {
	payload := &ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid HTTP response code: %s", resp.Status)
	}
	return fmt.Errorf("invalid HTTP response (%s): %s", resp.Status, payload.Description)
}
