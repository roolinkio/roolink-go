<div align="center">

# 🔗 Roolink Go SDK

### Enterprise-grade Akamai sensor generation for Go

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-purple?style=for-the-badge)](LICENSE)

**[📚 Documentation](https://docs.roolink.io/)** • **[💬 Discord](https://discord.gg/rooapi)**

---

</div>

## 🚀 Overview

Production-ready APIs for generating Akamai sensor data across web and mobile platforms. This SDK handles authentication, requests, and response parsing.

### 🎯 Supported Services

| Service | Features |
|---------|----------|
| **🌐 Web API** | Sensor generation, pixel challenges, SBSD, sec-cpt solving |
| **📱 BMP API** | iOS & Android mobile app sensor generation |

---

## 📦 Installation

```bash
go get github.com/roolinkio/roolink-go
```

## ⚡ Quick Start

```go
import roolink "github.com/roolinkio/roolink-go"

client := roolink.NewClient("your-api-key")
sensor, err := client.GenerateWebSensor(req)
```

---

## 🌐 Web API

```go
// Sensor Generation
client.GenerateWebSensor(req)

// Pixel Challenges
client.GeneratePixel(req)

// Sec-Cpt Challenges
client.SolveSecCpt(req)

// SBSD Challenges
client.SolveSBSD(req)

// Script Parsing
client.ParseScript(scriptContent)
```

---

## 📱 BMP API

```go
// iOS Sensors (android: false)
req := roolink.BMPSensorRequest{
    Android: false,
}
client.GenerateBMPSensor(req)

// Android Sensors (android: true)
req := roolink.BMPSensorRequest{
    Android: true,
}
client.GenerateBMPSensor(req)
```

---

## 🆘 Error Handling

```go
result, err := client.GenerateWebSensor(req)
if err != nil {
    // handle error
}
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

<div align="center">

**Built with 💜**

</div>
