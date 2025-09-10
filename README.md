# 🚀 Load Balancer in GO

<div align="center">

[Features](#-features) • [Quick Start](#-quick-start) • [Usage](#-usage) • [Configuration](#user-content-️-yaml-configuration-parameters) • [Strategies](#user-content--load-balancing-strategies) • [Health Checking](#-health-checking)

</div>

---

## 🎯 Overview

A **load balancer** built with Go that provides intelligent traffic distribution, automatic health checking, and multiple load balancing strategies.

## ✨ Features

- **Multiple Load Balancing Strategies:** Round Robin, Weighted, Smooth Weighted, Least Connections, Random
- **Intelligent Health Checking:** Concurrent checks, recovery, configurable intervals
- **Safe Architecture:** Thread-safe, graceful shutdown and logging
- **Flexible Configuration:** YAML, environment defaults, minimal setup

## 🚀 Quick Start

```bash
git clone https://github.com/franciscodelahoz/load-balancer.git
cd load-balancer
go mod tidy
go build -o load-balancer ./cmd/server
./load-balancer
```

## 📖 Usage

```bash
go run ./cmd/server/main.go
go run ./cmd/server/main.go -config=production.yaml
curl http://localhost:8080/
```

### Example Test Backends

```bash
# Terminal 1
python3 -m http.server 3001
# Terminal 2
python3 -m http.server 3002
```

---

## ⚙️ YAML Configuration Parameters

All configuration is done via a YAML file. Below, each parameter is explained in detail, including its purpose and default value if omitted.

---

### **server**

- **port**
  *(default: `8080`)*
  Port on which the load balancer HTTP server listens.

### **load_balancer**

- **strategy**
  *(default: `"round-robin"`)*
  Load balancing algorithm. Options: `"round-robin"`, `"weighted-round-robin"`, `"smooth-weighted-round-robin"`, `"least-connections"`, `"random"`.

### **backends**

- **url**
  *(required)*
  The URL of the backend service.

- **weight**
  *(default: `1`)*
  Relative weight for distributing traffic. Higher values mean more requests sent to this backend.

### **health_check**

- **enabled**
  *(default: `true`)*
  Enables or disables health checking.

- **interval**
  *(default: `10s`)*
  How often to perform health checks (Go duration format, e.g., `10s`, `1m`).

- **timeout**
  *(default: `5s`)*
  Timeout for each health check request.

- **path**
  *(default: `"/health"`)*
  Path to request on each backend for health checking.

- **method**
  *(default: `"GET"`)*
  HTTP method to use for health checks.

- **success_threshold**
  *(default: `3`)*
  Number of consecutive successful health checks required before a backend is marked healthy.

- **failure_threshold**
  *(default: `3`)*
  Number of consecutive failed health checks required before a backend is marked unhealthy.

---

### **Examples**

**Minimal Configuration:**

```yaml
backends:
  - url: "http://service-1:8080"
```

**Full Configuration:**

```yaml
server:
  port: 8080

load_balancer:
  strategy: "least-connections"

backends:
  - url: "http://service-1:8080"
    weight: 2

health_check:
  enabled: true
  interval: 30s
  timeout: 5s
  path: "/health"
  method: "GET"
  success_threshold: 5
  failure_threshold: 2
```

---

## ⚠️ Health Endpoint Guidance

The load balancer marks a backend as *healthy* when the configured health endpoint returns an HTTP 2xx status. If your application responds with 200 OK for unknown or invalid routes, the health check will always succeed and give a false positive.

### Recommendations:
- Expose a dedicated, lightweight health endpoint (e.g. /health) that returns 200 only when the service is actually healthy.
- Return 404/4xx for unknown or invalid paths.
- Keep health checks fast — avoid expensive operations.

### Quick test:

```bash
curl -i http://your-backend/health         # must return 200
curl -i http://your-backend/invalid-path   # must NOT return 200
```

**Go Example:**
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
})

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.NotFound(w, r)
})
```

---

## 🎯 Load Balancing Strategies

- `round-robin`
- `weighted-round-robin`
- `smooth-weighted-round-robin`
- `least-connections`
- `random`

Set the strategy in your YAML config:

```yaml
load_balancer:
  strategy: "smooth-weighted-round-robin"
```

---

## 🏥 Health Checking

**Config:**
```yaml
health_check:
  enabled: true
  interval: 10s
  timeout: 5s
  path: "/health"
  method: "GET"
```

**Best Practices:**
- Enable health checking for auto-recovery
- Use a dedicated health endpoint
- Set proper intervals and timeouts

---

## 📊 Monitoring & Metrics

- **Logging:**
  Logs show strategy, backend states, health results, routing decisions.

```
2025/09/10 10:50:24 🚀 Starting Load Balancer...
2025/09/10 10:50:24 ✅ Added backend: http://localhost:3002 (weight: 1)
2025/09/10 10:50:24 ✅ Registered backend for health checking: http://localhost:3002
2025/09/10 10:50:24 🏥 Health checker started with 1 backends
2025/09/10 10:50:24 🏥 Health checking enabled (interval: 10s)
2025/09/10 10:50:24 🚀 Load Balancer running on ::8080
2025/09/10 10:50:24 📊 Strategy: Smooth Weighted Round Robin
2025/09/10 10:50:24 🏢 Admin API: http://localhost::8080/admin/health
2025/09/10 10:50:29 ✅ Backend http://localhost:3002 health check passed (latency: 5.055478666s)
2025/09/10 10:51:39 ❌ Backend http://localhost:3002 health check failed: unexpected HTTP status from backend: 404
```

---

## 🏗️ Architecture

The diagram below shows the overall structure and flow:

```
                    Load Balancer
                         │
        ┌────────────────┼────────────────┐
        │                │                │
   Strategies        Health            Server
    Manager          Checker            Pool
        │                │                │
        └────────────────┼────────────────┘
                         │
                   ProxyHandler
                         │
                 ┌───────┼───────┐
                 │       │       │
            Backend1  Backend2  Backend3
```

- **Load Balancer:** Orchestrates incoming traffic and applies load balancing logic.
- **Strategies Manager:** Chooses the backend based on selected algorithm.
- **Health Checker:** Continuously checks backend health and availability.
- **Server Pool:** Maintains list and state of backend servers.
- **Proxy Handler:** Handles request forwarding and error responses.
- **Backends:** The actual application servers receiving requests.

---

## 🛠️ Development

**Project Structure:**

```
load-balancer/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── backend/                 # Backend management
│   ├── config/                  # Configuration handling
│   ├── handlers/                # HTTP handlers
│   ├── health/                  # Health checking
│   ├── loadbalancer/           # Core load balancer
│   └── strategies/             # Load balancing algorithms
├── config.yaml                 # Default configuration
└── README.md
```

**Build & Run:**
```bash
go run ./cmd/server/main.go
go build -o load-balancer ./cmd/server
./load-balancer
```

**Production Build:**
```bash
go build -ldflags="-w -s" -o load-balancer ./cmd/server
GOOS=linux GOARCH=amd64 go build -o load-balancer-linux ./cmd/server
```

**Roadmap:**
- [ ] Unit tests implementation
- [ ] Metrics endpoint (`/metrics`)
- [ ] Docker Compose setup and examples
- [ ] Performance benchmarks
- [ ] Admin API for runtime configuration
