package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type ImpervaPayloadRequest struct {
	Script string `json:"script"`
	// ScriptURL may be absolute or relative to PageURL.
	ScriptURL string `json:"scriptUrl"`
	PageURL   string `json:"pageUrl"`
	// UserAgent must be the one your requests to the site send, because the payload is built from a collected browser
	// with that exact User-Agent.
	UserAgent string `json:"userAgent"`
	// Language is the Accept-Language header your requests send, such as "en-US,en;q=0.9".
	Language string `json:"language"`
	// Pow is the site's base64 GPC challenge without the JSON quotes it arrives in. Leave it empty when the site has no
	// challenge.
	Pow string `json:"pow,omitempty"`
}

type ImpervaUTMVCRequest struct {
	Script string `json:"script"`
	// Cookies is the site's document.cookie. The API requires it even when it is empty.
	Cookies   string `json:"cookies"`
	UserAgent string `json:"userAgent"`
}

// ImpervaUTMVCResponse is the ___utmvc cookie, of which only Cookie or SetCookie is sent to the site. The fields
// after SetCookie describe how the value was built.
type ImpervaUTMVCResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	// Cookie is Name=Value, for a Cookie header.
	Cookie string `json:"cookie"`
	// SetCookie adds the cookie's expiry and path to Cookie, for a cookie jar.
	SetCookie      string `json:"setCookie"`
	Fingerprint    string `json:"fingerprint"`
	Digest         string `json:"digest"`
	Signature      string `json:"signature"`
	Salt           string `json:"salt"`
	Key            string `json:"key"`
	Version        string `json:"version"`
	WindowKeyCount int    `json:"windowKeyCount"`
}

// GenerateImpervaPayload returns the Reese84 submission as the API sent it, because it must be posted to the
// challenge endpoint unchanged.
func (c *Client) GenerateImpervaPayload(ctx context.Context, req ImpervaPayloadRequest) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/payload", DefaultImpervaBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint:errcheck // A failed close leaves nothing for the caller to handle.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) GenerateImpervaUTMVC(ctx context.Context, req ImpervaUTMVCRequest) (*ImpervaUTMVCResponse, error) {
	url := fmt.Sprintf("%s/api/v1/utmvc", DefaultImpervaBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result ImpervaUTMVCResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
