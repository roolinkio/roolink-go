package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// KasadaPayloadRequest asks for a Kasada /tl payload.
type KasadaPayloadRequest struct {
	Script    string `json:"script"`
	IPSURL    string `json:"ipsUrl"`
	Language  string `json:"language"`
	Referer   string `json:"referer"`
	UserAgent string `json:"userAgent"`
}

// KasadaPayloadResponse is a Kasada /tl request to send to the site.
type KasadaPayloadResponse struct {
	Headers map[string]string `json:"headers"`
	// Payload is the /tl request body. The API sends it as base64, which
	// encoding/json decodes into []byte.
	Payload []byte `json:"payload"`
}

// KasadaCDRequest asks for an x-kpsdk-cd header from the site's challenge
// values.
type KasadaCDRequest struct {
	ST     int64  `json:"st"`
	CT     string `json:"ct"`
	FC     string `json:"fc,omitempty"`
	Domain string `json:"domain"`
}

// KasadaCDResponse is an x-kpsdk-cd header value.
type KasadaCDResponse struct {
	// Payload is the x-kpsdk-cd header value, which is itself JSON, so send it
	// unchanged.
	Payload string `json:"payload"`
}

// KasadaBotIDRequest asks for an x-is-human header.
type KasadaBotIDRequest struct {
	Script string `json:"script"`
}

// KasadaBotIDResponse is sent as the x-is-human header by marshaling it with
// encoding/json. Its fields are in the API's order so that the marshaled header
// matches the API's response.
type KasadaBotIDResponse struct {
	B  float64 `json:"b"`
	V  float64 `json:"v"`
	E  string  `json:"e"`
	S  string  `json:"s"`
	D  float64 `json:"d"`
	VR string  `json:"vr"`
}

// GenerateKasadaPayload returns the Kasada /tl request for req.
func (c *Client) GenerateKasadaPayload(ctx context.Context, req KasadaPayloadRequest) (*KasadaPayloadResponse, error) {
	url := DefaultKasadaBaseURL + "/api/v1/payload"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaPayloadResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateKasadaCD returns the x-kpsdk-cd header for req.
func (c *Client) GenerateKasadaCD(ctx context.Context, req KasadaCDRequest) (*KasadaCDResponse, error) {
	url := DefaultKasadaBaseURL + "/api/v1/cd"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaCDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateKasadaBotID returns the x-is-human header for req.
func (c *Client) GenerateKasadaBotID(ctx context.Context, req KasadaBotIDRequest) (*KasadaBotIDResponse, error) {
	url := DefaultKasadaBaseURL + "/api/v1/botid"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaBotIDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
