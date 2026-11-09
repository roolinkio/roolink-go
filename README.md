# Roolink Go SDK

Go client for the Roolink Akamai Web (`web.roolink.io`), Akamai BMP (`bmp.roolink.io`) and Kasada (`ksd.roolink.io`) APIs.

Documentation: https://docs.roolink.io

## Install

```sh
go get github.com/roolinkio/roolink-go
```

## Usage

```go
package main

import (
	"context"
	"log"

	roolink "github.com/roolinkio/roolink-go"
)

func main() {
	ctx := context.Background()
	client := roolink.NewClient("your-api-key")

	scriptData, err := client.ParseScript(ctx, scriptBytes)
	if err != nil {
		log.Fatal(err)
	}

	sensor, err := client.GenerateWebSensor(ctx, roolink.WebSensorRequest{
		UserAgent:  userAgent,
		URL:        "https://example.com/",
		Abck:       abckCookie,
		BmSz:       bmSzCookie,
		ScriptData: scriptData,
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = sensor.Sensor
}
```

## Methods

| Method | Endpoint |
| --- | --- |
| `ParseScript(ctx, scriptBytes)` | `POST https://web.roolink.io/api/v1/parse` (gzip-compressed by default; disable with `WithCompression(false)`) |
| `GenerateWebSensor(ctx, WebSensorRequest)` | `POST https://web.roolink.io/api/v1/sensor` |
| `GeneratePixel(ctx, PixelRequest)` | `POST https://web.roolink.io/api/v1/pixel` |
| `SolveSecCpt(ctx, SecCptRequest)` | `POST https://web.roolink.io/api/v1/sec-cpt` |
| `SolveSBSD(ctx, SBSDRequest)` | `POST https://web.roolink.io/api/v1/sbsd` |
| `GenerateBMPSensor(ctx, BMPSensorRequest)` | `POST https://bmp.roolink.io/api/v1/sensor` (`Android: false` for iOS, `Android: true` for Android) |
| `GenerateKasadaPayload(ctx, KasadaPayloadRequest)` | `POST https://ksd.roolink.io/api/v1/payload` |
| `GenerateKasadaCD(ctx, KasadaCDRequest)` | `POST https://ksd.roolink.io/api/v1/cd` |
| `GenerateKasadaBotID(ctx, KasadaBotIDRequest)` | `POST https://ksd.roolink.io/api/v1/botid` |

Pass `WithHTTPClient(httpClient)` to override the default 30-second HTTP/2 client.

## Kasada

`GenerateKasadaPayload` decodes the payload from base64, so `Payload` is the `/tl` request body as is. Each protected request needs a new `GenerateKasadaCD`, whose `Payload` is the `x-kpsdk-cd` header value. For Vercel BotID, the `x-is-human` header is the response marshaled with `encoding/json`.

```go
answer, err := client.GenerateKasadaPayload(ctx, roolink.KasadaPayloadRequest{
	Script:    ipsScript,
	IPSURL:    ipsURL,
	Language:  "en-US",
	Referer:   "https://www.example.com/",
	UserAgent: userAgent,
})
if err != nil {
	log.Fatal(err)
}
// Post answer.Payload to /tl with answer.Headers, and keep its x-kpsdk-ct and x-kpsdk-st response headers.

cd, err := client.GenerateKasadaCD(ctx, roolink.KasadaCDRequest{ST: st, CT: ct, Domain: "www.example.com"})
if err != nil {
	log.Fatal(err)
}
_ = cd.Payload

botID, err := client.GenerateKasadaBotID(ctx, roolink.KasadaBotIDRequest{Script: botIDScript})
if err != nil {
	log.Fatal(err)
}
isHuman, err := json.Marshal(botID)
if err != nil {
	log.Fatal(err)
}
_ = string(isHuman)
```

## License

MIT — see [LICENSE](LICENSE).
