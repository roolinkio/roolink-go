package roolink

import (
	"bytes"
	"compress/gzip"
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

func TestCompression(t *testing.T) {
	const script = "var _ac=function(){};"

	// send parses the script with a client made with opts, and returns the request's header and body.
	send := func(opts ...ClientOption) (http.Header, []byte) {
		var header http.Header
		var body []byte
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			header = r.Header
			var err error
			if body, err = io.ReadAll(r.Body); err != nil {
				return nil, err
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
		})
		client := NewClient("roo_test", append(opts, WithHTTPClient(&http.Client{Transport: transport}))...)
		if _, err := client.ParseScript(context.Background(), []byte(script)); err != nil {
			t.Fatalf("ParseScript: %v", err)
		}
		return header, body
	}

	header, body := send()
	if got := header.Get("Content-Encoding"); got != "gzip" {
		t.Errorf("Content-Encoding = %q, want gzip", got)
	}
	zr, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("gzip.NewReader: %v", err)
	}
	if got, err := io.ReadAll(zr); err != nil || string(got) != script {
		t.Errorf("decompressed body = %q, %v; want %q", got, err, script)
	}

	header, body = send(WithCompression(false))
	if got := header.Get("Content-Encoding"); got != "" {
		t.Errorf("Content-Encoding = %q, want none", got)
	}
	if string(body) != script {
		t.Errorf("body = %q, want %q", body, script)
	}
}
