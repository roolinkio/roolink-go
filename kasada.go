package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type KasadaPayloadRequest struct {
	Script    string `json:"script"`
	IPSURL    string `json:"ipsUrl"`
	Language  string `json:"language"`
	Referer   string `json:"referer"`
	UserAgent string `json:"userAgent"`
}

type KasadaPayloadResponse struct {
	Headers map[string]string `json:"headers"`
	// Payload is the /tl request body. The API sends it as base64, which encoding/json decodes into []byte.
	Payload []byte `json:"payload"`
}

type KasadaCDRequest struct {
	ST     int64  `json:"st"`
	CT     string `json:"ct"`
	FC     string `json:"fc,omitempty"`
	Domain string `json:"domain"`
}

type KasadaCDResponse struct {
	// Payload is the x-kpsdk-cd header value, which is itself JSON, so send it unchanged.
	Payload string `json:"payload"`
}

type KasadaBotIDRequest struct {
	Script string `json:"script"`
}

// KasadaBotIDResponse is sent as the x-is-human header by marshaling it with encoding/json. Its fields are in the
// API's order so that the marshaled header matches the API's response.
type KasadaBotIDResponse struct {
	B  float64 `json:"b"`
	V  float64 `json:"v"`
	E  string  `json:"e"`
	S  string  `json:"s"`
	D  float64 `json:"d"`
	VR string  `json:"vr"`
}

func (c *Client) GenerateKasadaPayload(ctx context.Context, req KasadaPayloadRequest) (*KasadaPayloadResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payload", DefaultKasadaBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaPayloadResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GenerateKasadaCD(ctx context.Context, req KasadaCDRequest) (*KasadaCDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/cd", DefaultKasadaBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaCDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GenerateKasadaBotID(ctx context.Context, req KasadaBotIDRequest) (*KasadaBotIDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/botid", DefaultKasadaBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result KasadaBotIDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
