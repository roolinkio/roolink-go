package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// ScriptData is what ParseScript reads from an Akamai Web script, which
// WebSensorRequest needs.
type ScriptData struct {
	Ver string `json:"ver"`
	Key int64  `json:"key"`
	Dvc string `json:"dvc"`
	Din []int  `json:"din"`
}

// WebSensorRequest asks for an Akamai Web sensor.
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

// WebSensorResponse is an Akamai Web sensor to post to the site.
type WebSensorResponse struct {
	Sensor string `json:"sensor"`
}

// PixelRequest asks for an Akamai pixel sensor.
type PixelRequest struct {
	UserAgent            string `json:"userAgent"`
	Bazadebezolkohpepadr int    `json:"bazadebezolkohpepadr"`
	Hash                 string `json:"hash"`
}

// PixelResponse is an Akamai pixel sensor.
type PixelResponse struct {
	Sensor string `json:"sensor"`
}

// SecCptRequest is an Akamai sec-cpt challenge to solve.
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

// SecCptResponse is the solved sec-cpt challenge to post to the site.
type SecCptResponse struct {
	Token   string   `json:"token"`
	Answers []string `json:"answers"`
}

// SBSDRequest asks for an Akamai SBSD payload.
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

// SBSDResponse is the SBSD request body to post to the site.
type SBSDResponse struct {
	Body string `json:"body"`
}

// GenerateWebSensor returns an Akamai Web sensor for req.
func (c *Client) GenerateWebSensor(ctx context.Context, req WebSensorRequest) (*WebSensorResponse, error) {
	url := DefaultWebBaseURL + "/api/v1/sensor"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
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
	url := DefaultWebBaseURL + "/api/v1/pixel"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	defer closeBody(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return &PixelResponse{Sensor: string(body)}, nil
}

// SolveSecCpt solves the sec-cpt challenge in req.
func (c *Client) SolveSecCpt(ctx context.Context, req SecCptRequest) (*SecCptResponse, error) {
	url := DefaultWebBaseURL + "/api/v1/sec-cpt"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result SecCptResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SolveSBSD returns the SBSD payload for req.
func (c *Client) SolveSBSD(ctx context.Context, req SBSDRequest) (*SBSDResponse, error) {
	url := DefaultWebBaseURL + "/api/v1/sbsd"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result SBSDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ParseScript reads the values a sensor needs from an Akamai Web script.
func (c *Client) ParseScript(ctx context.Context, scriptContent []byte) (*ScriptData, error) {
	url := DefaultWebBaseURL + "/api/v1/parse"

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
