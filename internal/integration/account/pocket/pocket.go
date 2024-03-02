package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/integration/account"
	"github.com/ncarlier/readflow/pkg/mediatype"
)

type pocketProvider struct {
	key string
}

type pocketAuthorizeRequest struct {
	Key  string `json:"consumer_key"`
	Code string `json:"code"`
}

type pocketRequestResponse struct {
	Redirect string `json:"redirect"`
	Code     string `json:"code"`
}

func newPocketProvider(conf *config.IntegrationConfig) (account.Provider, error) {
	if conf.Pocket.ConsumerKey == "" {
		return nil, errors.New("pocket consumer key not set")
	}
	provider := &pocketProvider{
		key: conf.Pocket.ConsumerKey,
	}
	return provider, nil
}

// RequestHandler used for linking account request
func (p *pocketProvider) RequestHandler(w http.ResponseWriter, r *http.Request) error {
	redirect, ok := r.URL.Query()["redirect_uri"]
	if !ok || len(redirect[0]) < 1 {
		return errors.New("missing redirect_uri parameter")
	}
	params := url.Values{}
	params.Set("consumer_key", p.key)
	params.Set("redirect_uri", redirect[0])
	payload := strings.NewReader(params.Encode())
	resp, err := http.Post("https://getpocket.com/v3/oauth/request", mediatype.FormURLEncoded, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}
	code := values.Get("code")
	if code == "" {
		return errors.New("code is empty")
	}
	params.Del("consumer_key")
	params.Set("request_token", code)
	data := pocketRequestResponse{
		Redirect: "https://getpocket.com/auth/authorize?" + params.Encode(),
		Code:     code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(data)
	return nil
}

// AuthorizeHandler used for linking account authorization request
func (p *pocketProvider) AuthorizeHandler(w http.ResponseWriter, r *http.Request) error {
	code, ok := r.URL.Query()["code"]
	if !ok || len(code[0]) < 1 {
		return errors.New("missing code parameter")
	}
	params := pocketAuthorizeRequest{
		Key:  p.key,
		Code: code[0],
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(params)

	req, err := http.NewRequest("POST", "https://getpocket.com/v3/oauth/authorize", b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mediatype.JSON)
	req.Header.Set("X-Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Forward response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	return nil
}

func init() {
	account.Register("pocket", &account.Def{
		Name:   "Pocket",
		Desc:   "Put knowledge in your Pocket.",
		Create: newPocketProvider,
	})
}
