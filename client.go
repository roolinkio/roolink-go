// Package roolink is the Go client for the Roolink APIs, which solve Akamai Web,
// Akamai BMP, Kasada, and Imperva bot challenges.
package roolink

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultTimeout bounds each request of a client built without
// WithHTTPClient. The base URLs are the APIs' production endpoints.
const (
	DefaultTimeout        = 30 * time.Second
	DefaultBMPBaseURL     = "https://bmp.roolink.io"
	DefaultWebBaseURL     = "https://web.roolink.io"
	DefaultKasadaBaseURL  = "https://ksd.roolink.io"
	DefaultImpervaBaseURL = "https://imperva.roolink.io"
)

// Client calls the Roolink APIs with one API key. It is safe for concurrent
// use.
type Client struct {
	apiKey     string
	httpClient *http.Client
	compress   bool
}

// ClientOption configures a Client in NewClient.
type ClientOption func(*Client)

// WithHTTPClient makes the client send requests with httpClient instead of its
// own.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithCompression sets whether ParseScript gzips the script payload (enabled by
// default).
func WithCompression(compress bool) ClientOption {
	return func(c *Client) {
		c.compress = compress
	}
}

// NewClient returns a Client that authenticates with apiKey.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	transport := &http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig:   &tls.Config{},
	}

	client := &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		},
		compress: true,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) doRequest(ctx context.Context, url string, body io.Reader, contentType string, compress bool) (*http.Response, error) {
	if compress {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		if _, err := io.Copy(zw, body); err != nil {
			return nil, fmt.Errorf("compress request: %w", err)
		}
		if err := zw.Close(); err != nil {
			return nil, fmt.Errorf("compress request: %w", err)
		}
		body = &buf
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "roolink-go/"+Version)
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	return resp, nil
}

func parseResponse(resp *http.Response, target any) error {
	defer closeBody(resp.Body)

	if resp.StatusCode >= 400 {
		// The body only adds detail to the status, so a failed read omits it.
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// closeBody closes a response body the caller is done with. It ignores the
// error, because a failed close leaves nothing for the caller to handle.
func closeBody(body io.Closer) {
	_ = body.Close()
}
