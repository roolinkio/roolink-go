package roolink

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGenerateImpervaPayload(t *testing.T) {
	// Callers post the submission to Imperva unchanged, so the client must return the bytes the API sent.
	const submission = `{"solution":{"interrogation":{"p":"AAEC/w==","st":1788893087,"sr":3420118839,"cr":2045817226,"og":2},"version":"stable"},"old_token":null,"error":null,"errors":[],"performance":{"1030820235":41,"4203433130":12,"total":43}}`

	for _, tc := range []struct {
		name string
		req  ImpervaPayloadRequest
		want string
	}{
		{
			name: "with pow",
			req: ImpervaPayloadRequest{
				Script:    "(function(){})();",
				ScriptURL: "/Hedge-the-Stage-Thy-Seed/script.js",
				PageURL:   "https://www.example.com/",
				UserAgent: "Mozilla/5.0",
				Language:  "en-US,en;q=0.9",
				Pow:       "eyJjIjoiYWJjIiwiZCI6NH0=",
			},
			want: `{"script":"(function(){})();","scriptUrl":"/Hedge-the-Stage-Thy-Seed/script.js","pageUrl":"https://www.example.com/","userAgent":"Mozilla/5.0","language":"en-US,en;q=0.9","pow":"eyJjIjoiYWJjIiwiZCI6NH0="}`,
		},
		{
			name: "without pow",
			req: ImpervaPayloadRequest{
				Script:    "(function(){})();",
				ScriptURL: "/Hedge-the-Stage-Thy-Seed/script.js",
				PageURL:   "https://www.example.com/",
				UserAgent: "Mozilla/5.0",
				Language:  "en-US,en;q=0.9",
			},
			want: `{"script":"(function(){})();","scriptUrl":"/Hedge-the-Stage-Thy-Seed/script.js","pageUrl":"https://www.example.com/","userAgent":"Mozilla/5.0","language":"en-US,en;q=0.9"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := jsonClient(t, "https://imperva.roolink.io/api/v1/payload", tc.want, submission).GenerateImpervaPayload(context.Background(), tc.req)
			if err != nil {
				t.Fatalf("GenerateImpervaPayload: %v", err)
			}
			if string(body) != submission {
				t.Errorf("body = %s, want %s", body, submission)
			}
		})
	}
}

func TestGenerateImpervaPayloadError(t *testing.T) {
	// An error body is JSON as well, and posting it to Imperva as a submission would fail without saying why.
	const reply = `{"code":"device_not_found","error":"no collected browser has this User-Agent"}`
	client := NewClient("roo_test", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader(reply)),
		}, nil
	})}))

	body, err := client.GenerateImpervaPayload(context.Background(), ImpervaPayloadRequest{})
	if err == nil {
		t.Fatalf("GenerateImpervaPayload = %s, want an error", body)
	}
	if !strings.Contains(err.Error(), "device_not_found") {
		t.Errorf("error = %v, want it to include the API's code", err)
	}
}

func TestGenerateImpervaUTMVC(t *testing.T) {
	const reply = `{"name":"___utmvc","value":"jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==","cookie":"___utmvc=jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==","setCookie":"___utmvc=jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==; expires=Mon, 09 Nov 2026 19:04:40 GMT; path=/","fingerprint":"navigator%3Dtrue,navigator.vendor%3DGoogle%20Inc.","digest":"7261","signature":"9da49ba49f5fa992a3a663a1a3a7a9","salt":"fresh-salt-plus","key":"fresh","version":"version-one","windowKeyCount":289}` + "\n"
	want := ImpervaUTMVCResponse{
		Name:           "___utmvc",
		Value:          "jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==",
		Cookie:         "___utmvc=jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==",
		SetCookie:      "___utmvc=jxKgPCxkaWdlc3Q9NzI2MSxzPTlkYTQ5YmE0OWY1ZmE5OTJhM2E2NjNhMWEzYTdhOQ==; expires=Mon, 09 Nov 2026 19:04:40 GMT; path=/",
		Fingerprint:    "navigator%3Dtrue,navigator.vendor%3DGoogle%20Inc.",
		Digest:         "7261",
		Signature:      "9da49ba49f5fa992a3a663a1a3a7a9",
		Salt:           "fresh-salt-plus",
		Key:            "fresh",
		Version:        "version-one",
		WindowKeyCount: 289,
	}

	for _, tc := range []struct {
		name string
		req  ImpervaUTMVCRequest
		want string
	}{
		{
			name: "with cookies",
			req:  ImpervaUTMVCRequest{Script: "(function(){})();", Cookies: "incap_ses_1234_567890=abc; visid_incap_567890=def", UserAgent: "Mozilla/5.0"},
			want: `{"script":"(function(){})();","cookies":"incap_ses_1234_567890=abc; visid_incap_567890=def","userAgent":"Mozilla/5.0"}`,
		},
		{
			// The API requires cookies, so an empty document.cookie must still be sent.
			name: "without cookies",
			req:  ImpervaUTMVCRequest{Script: "(function(){})();", UserAgent: "Mozilla/5.0"},
			want: `{"script":"(function(){})();","cookies":"","userAgent":"Mozilla/5.0"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := jsonClient(t, "https://imperva.roolink.io/api/v1/utmvc", tc.want, reply).GenerateImpervaUTMVC(context.Background(), tc.req)
			if err != nil {
				t.Fatalf("GenerateImpervaUTMVC: %v", err)
			}
			if *resp != want {
				t.Errorf("response = %+v, want %+v", *resp, want)
			}
		})
	}
}
