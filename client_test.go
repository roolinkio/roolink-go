package roolink

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// jsonClient returns a client whose one request must POST want as JSON to url, and which answers reply.
func jsonClient(t *testing.T, url, want, reply string) *Client {
	t.Helper()
	return NewClient("roo_test", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.String(); got != url {
			t.Errorf("URL = %q, want %q", got, url)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("x-api-key"); got != "roo_test" {
			t.Errorf("x-api-key = %q, want roo_test", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		var gotJSON, wantJSON any
		if err := json.Unmarshal(body, &gotJSON); err != nil {
			t.Errorf("body %q is not JSON: %v", body, err)
		}
		if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
			t.Fatalf("want %q is not JSON: %v", want, err)
		}
		if !reflect.DeepEqual(gotJSON, wantJSON) {
			t.Errorf("body = %s, want %s", body, want)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader(reply)),
		}, nil
	})}))
}

func TestUserAgent(t *testing.T) {
	var got string
	client := NewClient("roo_test", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		got = r.Header.Get("User-Agent")
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}))

	if _, err := client.ParseScript(context.Background(), []byte("var _ac=function(){};")); err != nil {
		t.Fatalf("ParseScript: %v", err)
	}
	if want := "roolink-go/" + Version; got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
}
