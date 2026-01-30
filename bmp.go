package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// BMPSensorRequest represents a request to generate an Akamai BMP sensor
type BMPSensorRequest struct {
	AppName  string `json:"app"`
	Proxy    string `json:"proxy"`
	Language string `json:"language,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	Android  bool   `json:"android,omitempty"`
	IPad     bool   `json:"ipad,omitempty"`
}

// Cookie represents a cookie returned from BMP sensor generation
type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

// BMPSensorResponse represents the response from BMP sensor generation
type BMPSensorResponse struct {
	Sensor   string `json:"sensor"`
	Platform string `json:"platform"`

	// iOS specific fields
	IOS             string   `json:"ios,omitempty"`
	KernelOsRelease string   `json:"kernelOsRelease,omitempty"`
	KernelOsVersion string   `json:"kernelOsVersion,omitempty"`
	MachineID       string   `json:"machineId,omitempty"`
	Cookies         []Cookie `json:"cookies,omitempty"`

	// Android specific fields
	Android            string `json:"android,omitempty"`
	DeviceModel        string `json:"deviceModel,omitempty"`
	DeviceManufacturer string `json:"deviceManufacturer,omitempty"`
	SDKVersion         string `json:"sdkVersion,omitempty"`

	// Common fields
	DeviceID       string `json:"deviceId"`
	AppVersion     string `json:"appVersion"`
	AppVersionCode int 	  `json:"appVersionCode,omitempty"`
	ScreenHeight   int    `json:"screenHeight"`
	ScreenWidth    int    `json:"screenWidth"`
	Language       string `json:"language"`
}

// GenerateBMPSensor generates an Akamai BMP sensor for mobile apps
func (c *Client) GenerateBMPSensor(ctx context.Context, req BMPSensorRequest) (*BMPSensorResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sensor", DefaultBMPBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", url, bytes.NewBuffer(jsonData), "application/json")
	if err != nil {
		return nil, err
	}

	var result BMPSensorResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
