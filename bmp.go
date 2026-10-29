package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type BMPSensorRequest struct {
	AppName  string `json:"app"`
	Proxy    string `json:"proxy"`
	Language string `json:"language,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	Android  bool   `json:"android,omitempty"`
	IPad     bool   `json:"ipad,omitempty"`
}

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

type BMPSensorResponse struct {
	Sensor   string `json:"sensor"`
	Platform string `json:"platform"`

	IOS             string   `json:"ios,omitempty"`
	KernelOsRelease string   `json:"kernelOsRelease,omitempty"`
	KernelOsVersion string   `json:"kernelOsVersion,omitempty"`
	MachineID       string   `json:"machineId,omitempty"`
	Cookies         []Cookie `json:"cookies,omitempty"`

	Android            string `json:"android,omitempty"`
	DeviceModel        string `json:"deviceModel,omitempty"`
	DeviceManufacturer string `json:"deviceManufacturer,omitempty"`
	SDKVersion         string `json:"sdkVersion,omitempty"`

	DeviceID       string `json:"deviceId"`
	AppVersion     string `json:"appVersion"`
	AppVersionCode int    `json:"appVersionCode,omitempty"`
	ScreenHeight   int    `json:"screenHeight"`
	ScreenWidth    int    `json:"screenWidth"`
	Language       string `json:"language"`
}

func (c *Client) GenerateBMPSensor(ctx context.Context, req BMPSensorRequest) (*BMPSensorResponse, error) {
	url := fmt.Sprintf("%s/api/v1/sensor", DefaultBMPBaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(jsonData), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result BMPSensorResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
