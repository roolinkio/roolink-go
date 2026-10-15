package roolink

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestGeneratePixelReturnsTextBody(t *testing.T) {
	const body = "ap=true&bt=0&fonts=1&fh=abc"
	client := NewClient("roo_test", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"text/plain"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}))

	resp, err := client.GeneratePixel(context.Background(), PixelRequest{})
	if err != nil {
		t.Fatalf("GeneratePixel: %v", err)
	}
	if resp.Sensor != body {
		t.Fatalf("Sensor = %q, want %q", resp.Sensor, body)
	}
}
