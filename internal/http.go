package internal

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
)

type HTTPSession interface {
	Fetch(ctx context.Context, method, url, contentType string, body []byte, authorization string) (response []byte, err error)
	FetchBot(ctx *EventContext, method, url, contentType string, body []byte) (response []byte, err error)

	FetchJSON(ctx context.Context, method, url string, payload interface{}, structure interface{}, authorization string) (err error)
	FetchJSONBot(ctx *EventContext, method, url string, payload interface{}, structure interface{}) (err error)
}

type TwilightProxy struct {
	HTTP       *http.Client
	APIVersion string
	URLHost    string
	URLScheme  string
	UserAgent  string
}

func NewTwilightProxy(url url.URL) (httpSession HTTPSession) {
	return &TwilightProxy{
		HTTP: &http.Client{
			Timeout: 20 * time.Second,
		},
		APIVersion: "9",
		URLHost:    url.Host,
		URLScheme:  url.Scheme,
		UserAgent:  "Sandwich/" + VERSION + " (github.com/WelcomerTeam/Sandwich)",
	}
}

// Fetch sends a request to the TwilightProxy and returns the raw response. For most requests,
// you will want to use FetchJSON.
func (tl *TwilightProxy) Fetch(ctx context.Context, method, url, contentType string, body []byte, authorization string) (response []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, xerrors.Errorf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", tl.UserAgent)

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := tl.HTTP.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("Failed to do request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidToken
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("Failed to read body: %v", err)
	}

	return response, err
}

// FetchJSON works similar to Fetch however will turn the payload into JSON (if passed) and
// converts the output into the structure passed.
func (tl *TwilightProxy) FetchJSON(ctx context.Context, method, url string, payload interface{}, structure interface{}, authorization string) (err error) {
	var body []byte

	if payload != nil {
		body, err = jsoniter.Marshal(payload)
		if err != nil {
			return xerrors.Errorf("Failed to marshal payload: %v", err)
		}
	}

	response, err := tl.Fetch(ctx, method, url, "application/json", body, authorization)
	if err != nil {
		return xerrors.Errorf("Failed to fetch: %v", err)
	}

	err = jsoniter.Unmarshal(response, &structure)
	if err != nil {
		return xerrors.Errorf("Failed to unmarshal response: %v", err)
	}

	return nil
}

// FetchBot is similar to Fetch() however automatically passes in the proper token.
func (tl *TwilightProxy) FetchBot(ctx *EventContext, method, url, contentType string, body []byte) (response []byte, err error) {
	return tl.Fetch(ctx.Context, method, url, contentType, body, "Bot "+ctx.Identifier.Token)
}

// FetchBot is similar to FetchJSON() however automatically passes in the proper token.
func (tl *TwilightProxy) FetchJSONBot(ctx *EventContext, method, url string, payload interface{}, structure interface{}) (err error) {
	return tl.FetchJSON(ctx.Context, method, url, payload, structure, "Bot "+ctx.Identifier.Token)
}
