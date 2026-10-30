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

const (
	DefaultTimeout    = 30 * time.Second
	DefaultBMPBaseURL = "https://bmp.roolink.io"
	DefaultWebBaseURL = "https://web.roolink.io"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	compress   bool
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithCompression sets whether ParseScript gzips the script payload (enabled by default).
func WithCompression(compress bool) ClientOption {
	return func(c *Client) {
		c.compress = compress
	}
}

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
			return nil, fmt.Errorf("failed to compress request: %w", err)
		}
		if err := zw.Close(); err != nil {
			return nil, fmt.Errorf("failed to compress request: %w", err)
		}
		body = &buf
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "roolink-go/"+Version)
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close() //nolint:errcheck // A failed close leaves nothing for the caller to handle.

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
