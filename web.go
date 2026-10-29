package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type ScriptData struct {
	Ver string `json:"ver"`
	Key int64  `json:"key"`
	Dvc string `json:"dvc"`
	Din []int  `json:"din"`
}

type WebSensorRequest struct {
	UserAgent  string      `json:"userAgent"`
	URL        string      `json:"url"`
	Abck       string      `json:"_abck"`
	BmSz       string      `json:"bm_sz"`
	ScriptData *ScriptData `json:"scriptData"`
	SecCpt     bool        `json:"sec_cpt,omitempty"`
	Index      int         `json:"index,omitempty"`
	Stepper    bool        `json:"stepper,omitempty"`
	Keyboard   bool        `json:"keyboard,omitempty"`
	Count      bool        `json:"count,omitempty"`
	Language   string      `json:"language,omitempty"`
	Flags      string      `json:"flags,omitempty"`
	ScriptURL  string      `json:"scriptUrl,omitempty"`
}

type WebSensorResponse struct {
	Sensor string `json:"sensor"`
}

type PixelRequest struct {
	UserAgent            string `json:"userAgent"`
	Bazadebezolkohpepadr int    `json:"bazadebezolkohpepadr"`
	Hash                 string `json:"hash"`
}

type PixelResponse struct {
	Sensor string `json:"sensor"`
}

type SecCptRequest struct {
	SecCpChallenge string `json:"sec-cp-challenge,omitempty"`
	Provider       string `json:"provider,omitempty"`
	BrandingURL    string `json:"branding_url_content,omitempty"`
	ChlgDuration   int    `json:"chlg_duration,omitempty"`
	Token          string `json:"token"`
	Timestamp      int    `json:"timestamp"`
	Nonce          string `json:"nonce"`
	Difficulty     int64  `json:"difficulty"`
	Count          int    `json:"count,omitempty"`
	Timeout        int64  `json:"timeout,omitempty"`
	CPU            bool   `json:"cpu,omitempty"`
	Cookie         string `json:"cookie"`
}

type SecCptResponse struct {
	Token   string   `json:"token"`
	Answers []string `json:"answers"`
}

type SBSDRequest struct {
	Vid        string `json:"vid"`
	UserAgent  string `json:"userAgent"`
	BmO        string `json:"bm_o"`
	Legacy     bool   `json:"legacy,omitempty"`
	Url        string `json:"url"`
	ScriptHash string `json:"script_hash,omitempty"`
	ScriptUrl  string `json:"script_url,omitempty"`
	Language   string `json:"language,omitempty"`
}

type SBSDResponse struct {
	Body string `json:"body"`
}

func (c *Client) GenerateWebSensor(ctx context.Context, req WebSensorRequest) (*WebSensorResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sensor", DefaultWebBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result WebSensorResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GeneratePixel posts a pixel request and reads the text/plain sensor response.
func (c *Client) GeneratePixel(ctx context.Context, req PixelRequest) (*PixelResponse, error) {
	url := fmt.Sprintf("%s/api/v1/pixel", DefaultWebBaseURL)

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

	return &PixelResponse{Sensor: string(body)}, nil
}

func (c *Client) SolveSecCpt(ctx context.Context, req SecCptRequest) (*SecCptResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sec-cpt", DefaultWebBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result SecCptResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) SolveSBSD(ctx context.Context, req SBSDRequest) (*SBSDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sbsd", DefaultWebBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result SBSDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ParseScript(ctx context.Context, scriptContent []byte) (*ScriptData, error) {
	url := fmt.Sprintf("%s/api/v1/parse", DefaultWebBaseURL)

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(scriptContent), "text/plain", c.compress)
	if err != nil {
		return nil, err
	}

	var result ScriptData
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
