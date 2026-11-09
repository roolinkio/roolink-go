package roolink

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestGenerateKasadaPayload(t *testing.T) {
	client := jsonClient(t, "https://ksd.roolink.io/api/v1/payload",
		`{"script":"(function(){})();","ipsUrl":"https://www.example.com/a/b/ips.js?x-kpsdk-v=j-1.1.0","language":"en-US","referer":"https://www.example.com/","userAgent":"Mozilla/5.0"}`,
		`{"headers":{"x-kpsdk-ct":"ct","x-kpsdk-v":"j-1.1.0"},"payload":"APv/EA=="}`+"\n",
	)

	resp, err := client.GenerateKasadaPayload(context.Background(), KasadaPayloadRequest{
		Script:    "(function(){})();",
		IPSURL:    "https://www.example.com/a/b/ips.js?x-kpsdk-v=j-1.1.0",
		Language:  "en-US",
		Referer:   "https://www.example.com/",
		UserAgent: "Mozilla/5.0",
	})
	if err != nil {
		t.Fatalf("GenerateKasadaPayload: %v", err)
	}
	if want := map[string]string{"x-kpsdk-ct": "ct", "x-kpsdk-v": "j-1.1.0"}; !reflect.DeepEqual(resp.Headers, want) {
		t.Errorf("Headers = %v, want %v", resp.Headers, want)
	}
	if want := []byte{0x00, 0xfb, 0xff, 0x10}; !bytes.Equal(resp.Payload, want) {
		t.Errorf("Payload = %x, want %x", resp.Payload, want)
	}
}

func TestGenerateKasadaCD(t *testing.T) {
	const cd = `{"workTime":1788893087900,"id":"0123456789abcdef0123456789abcdef","answers":[3,7],"duration":2,"d":120,"st":1788893087778,"rst":1788893087898}`
	reply, err := json.Marshal(map[string]string{"payload": cd})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		req  KasadaCDRequest
		want string
	}{
		{
			name: "with fc",
			req:  KasadaCDRequest{ST: 1788893087778, CT: "0aGzYc9d3d1b8f6f4e0e2d7c1b9a5f3e", FC: "eyJ4IjoxfQ==", Domain: "www.example.com"},
			want: `{"st":1788893087778,"ct":"0aGzYc9d3d1b8f6f4e0e2d7c1b9a5f3e","fc":"eyJ4IjoxfQ==","domain":"www.example.com"}`,
		},
		{
			name: "without fc",
			req:  KasadaCDRequest{ST: 1788893087778, CT: "0aGzYc9d3d1b8f6f4e0e2d7c1b9a5f3e", Domain: "www.example.com"},
			want: `{"st":1788893087778,"ct":"0aGzYc9d3d1b8f6f4e0e2d7c1b9a5f3e","domain":"www.example.com"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := jsonClient(t, "https://ksd.roolink.io/api/v1/cd", tc.want, string(reply)).GenerateKasadaCD(context.Background(), tc.req)
			if err != nil {
				t.Fatalf("GenerateKasadaCD: %v", err)
			}
			if resp.Payload != cd {
				t.Errorf("Payload = %q, want %q", resp.Payload, cd)
			}
		})
	}
}

func TestGenerateKasadaBotID(t *testing.T) {
	const header = `{"b":0,"v":0.08084381062192014,"e":"eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..uydqejFMULzuEHyK.Eh2s.aM1k","s":"q+/0Zm9v","d":2,"vr":"3"}`
	client := jsonClient(t, "https://ksd.roolink.io/api/v1/botid", `{"script":"(function(){})();"}`, header+"\n")

	resp, err := client.GenerateKasadaBotID(context.Background(), KasadaBotIDRequest{Script: "(function(){})();"})
	if err != nil {
		t.Fatalf("GenerateKasadaBotID: %v", err)
	}
	want := KasadaBotIDResponse{
		B:  0,
		V:  0.08084381062192014,
		E:  "eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..uydqejFMULzuEHyK.Eh2s.aM1k",
		S:  "q+/0Zm9v",
		D:  2,
		VR: "3",
	}
	if *resp != want {
		t.Errorf("response = %+v, want %+v", *resp, want)
	}

	// The x-is-human header is the response marshaled back to JSON, which must match what the API sent.
	got, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(got) != header {
		t.Errorf("marshaled = %s, want %s", got, header)
	}
}
