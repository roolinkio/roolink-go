package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// BMPSensorRequest asks for an Akamai BMP sensor for a mobile app.
type BMPSensorRequest struct {
	AppName  string `json:"app"`
	Proxy    string `json:"proxy"`
	Language string `json:"language,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	Android  bool   `json:"android,omitempty"`
	IPad     bool   `json:"ipad,omitempty"`
}

// Cookie is a cookie the app must set before it sends the sensor.
type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

// BMPSensorResponse is a BMP sensor and the device it was built for. The iOS
// and Android fields are set only for that platform.
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

// GenerateBMPSensor returns an Akamai BMP sensor for req.
func (c *Client) GenerateBMPSensor(ctx context.Context, req BMPSensorRequest) (*BMPSensorResponse, error) {
	url := DefaultBMPBaseURL + "/api/v1/sensor"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, url, bytes.NewBuffer(payload), "application/json", false)
	if err != nil {
		return nil, err
	}

	var result BMPSensorResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
