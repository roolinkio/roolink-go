package roolink

import (
	"fmt"
)

// ScriptData represents parsed Akamai script configuration data
type ScriptData struct {
	Ver string `json:"ver"`
	Key int64  `json:"key"`
	Dvc string `json:"dvc"`
	Din []int  `json:"din"`
}

// WebSensorRequest represents a request to generate an Akamai web sensor
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

// WebSensorResponse represents the response from web sensor generation
type WebSensorResponse struct {
	Sensor string `json:"sensor"`
}

// PixelRequest represents a request to generate pixel data
type PixelRequest struct {
	UserAgent            string `json:"userAgent"`
	Bazadebezolkohpepadr int    `json:"bazadebezolkohpepadr"`
	Hash                 string `json:"hash"`
}

// PixelResponse represents the response from pixel generation
type PixelResponse struct {
	Sensor string `json:"sensor"`
}

// SecCptRequest represents a request to solve a sec-cpt challenge
type SecCptRequest struct {
	SecCpChallenge string `json:"sec-cp-challenge,omitempty"`
	Provider       string `json:"provider,omitempty"`
	BrandingURL    string `json:"branding_url_content,omitempty"`
	ChlgDuration   int    `json:"chlg_duration,omitempty"`
	Token          string `json:"token"`
	Timestamp      int    `json:"timestamp"`
	Nonce          string `json:"nonce"`
	Difficulty     int64  `json:"difficulty"`
	Timeout        int64  `json:"timeout,omitempty"`
	CPU            bool   `json:"cpu,omitempty"`
	Cookie         string `json:"cookie"`
}

// SecCptResponse represents the response from sec-cpt challenge
type SecCptResponse struct {
	Token   string   `json:"token"`
	Answers []string `json:"answers"`
}

// SBSDRequest represents a request to solve SBSD challenge
type SBSDRequest struct {
	Vid        string `json:"vid"`
	UserAgent  string `json:"userAgent"`
	BmO        string `json:"bm_o"`
	Legacy     bool   `json:"legacy,omitempty"`
	Url        string `json:"url"`
	ScriptHash string `json:"script_hash,omitempty"`
	ScriptUrl  string `json:"script_url,omitempty"`
}

// SBSDResponse represents the response from SBSD challenge
type SBSDResponse struct {
	Body string `json:"body"`
}

// ParseRequest represents a request to parse Akamai script
type ParseRequest struct {
	ScriptContent string `json:"scriptContent"`
}

// ParseResponse represents the response from script parsing
type ParseResponse struct {
	ScriptData ScriptData `json:"scriptData"`
}

// GenerateWebSensor generates an Akamai web sensor
func (c *Client) GenerateWebSensor(req WebSensorRequest) (*WebSensorResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sensor", DefaultWebBaseURL)

	resp, err := c.doRequest("POST", url, req)
	if err != nil {
		return nil, err
	}

	var result WebSensorResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GeneratePixel generates pixel sensor data
func (c *Client) GeneratePixel(req PixelRequest) (*PixelResponse, error) {
	url := fmt.Sprintf("%s/api/v1/pixel", DefaultWebBaseURL)

	resp, err := c.doRequest("POST", url, req)
	if err != nil {
		return nil, err
	}

	var result PixelResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SolveSecCpt solves a sec-cpt crypto challenge
func (c *Client) SolveSecCpt(req SecCptRequest) (*SecCptResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sec-cpt", DefaultWebBaseURL)

	resp, err := c.doRequest("POST", url, req)
	if err != nil {
		return nil, err
	}

	var result SecCptResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SolveSBSD solves an SBSD challenge
func (c *Client) SolveSBSD(req SBSDRequest) (*SBSDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sbsd", DefaultWebBaseURL)

	resp, err := c.doRequest("POST", url, req)
	if err != nil {
		return nil, err
	}

	var result SBSDResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ParseScript parses an Akamai script and returns script data
func (c *Client) ParseScript(scriptContent []byte) (*ParseResponse, error) {
	url := fmt.Sprintf("%s/api/v1/parse", DefaultWebBaseURL)

	req := map[string]string{
		"scriptContent": string(scriptContent),
	}

	resp, err := c.doRequest("POST", url, req)
	if err != nil {
		return nil, err
	}

	var result ParseResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
