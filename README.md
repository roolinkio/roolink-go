# Roolink Go SDK

Go client for the Roolink Akamai Web (`web.roolink.io`) and Akamai BMP (`bmp.roolink.io`) APIs.

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

Pass `WithHTTPClient(httpClient)` to override the default 30-second HTTP/2 client.

## License

MIT — see [LICENSE](LICENSE).
