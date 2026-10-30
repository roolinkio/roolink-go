package roolink

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

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
