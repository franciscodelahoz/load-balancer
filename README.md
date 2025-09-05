# 🚀 Load Balancer in GO

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)


[Features](#-features) • [Quick Start](#-quick-start) • [Configuration](#-configuration) • [Strategies](#-load-balancing-strategies) • [Health Checking](#-health-checking)

</div>

---

## 🎯 Overview

A **load balancer** built with Go that provides intelligent traffic distribution, automatic health checking, and multiple load balancing strategies.

## ✨ Features

- **🎯 Multiple Load Balancing Strategies**
  - Round Robin
  - Weighted Round Robin
  - Smooth Weighted Round Robin
  - Least Connections
  - Random

- **🏥 Intelligent Health Checking**
  - Configurable health check intervals
  - Automatic failure detection
  - Backend recovery monitoring
  - Concurrent health checks

- **⚙️ Safe Architecture**
  - Thread-safe operations
  - Graceful shutdown
  - Comprehensive error handling
  - Request metrics and logging

- **📝 Flexible Configuration**
  - YAML-based configuration
  - Environment-based defaults
  - Validation and warnings
  - Minimal configuration required

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Network access to backend servers

### Installation

```bash
# Clone the repository
git clone https://github.com/franciscodelahoz/load-balancer.git
cd load-balancer

# Install dependencies
go mod tidy

# Build the load balancer
go build -o load-balancer ./cmd/server

# Run with default configuration
./load-balancer
```

### Docker (Optional)

```bash
# Build Docker image
docker build -t franciscodelahoz/load-balancer .

# Run container
docker run -p 8080:8080 -v $(pwd)/config.yaml:/config.yaml franciscodelahoz/load-balancer
```

## 📖 Usage

### Basic Usage

```bash
# Run with default configuration
go run ./cmd/server/main.go

# Run with custom config file
go run ./cmd/server/main.go -config=production.yaml

# Test the load balancer
curl http://localhost:8080/
```

### Example Backend Servers

Start some test backends:

```bash
# Terminal 1 - Backend 1
python3 -m http.server 3001

# Terminal 2 - Backend 2
python3 -m http.server 3002

# Terminal 3 - Backend 3
python3 -m http.server 3003
```

## ⚙️ Configuration

### Basic Configuration

```yaml
# config.yaml
server:
  port: 8080

load_balancer:
  strategy: "smooth-weighted-round-robin"

backends:
  - url: "http://localhost:3001"
    weight: 1
  - url: "http://localhost:3002"
    weight: 2
  - url: "http://localhost:3003"
    weight: 3

health_check:
  enabled: true
  interval: 10s
  timeout: 5s
  path: "/"
  method: "GET"
```

### Production Configuration

```yaml
# production.yaml
server:
  port: 80

load_balancer:
  strategy: "least-connections"

backends:
  - url: "https://api-1.company.com"
    weight: 3
  - url: "https://api-2.company.com"
    weight: 2
  - url: "https://api-3.company.com"
    weight: 1

health_check:
  enabled: true
  interval: 30s
  timeout: 10s
  path: "/health"
  method: "GET"
```

### Minimal Configuration

```yaml
# Only specify what you need - rest uses defaults
backends:
  - url: "http://service-1:8080"
  - url: "http://service-2:8080"
```

## 🎯 Load Balancing Strategies

### Round Robin
Distributes requests sequentially across backends.

```yaml
load_balancer:
  strategy: "round-robin"
```

### Weighted Round Robin
Distributes requests based on backend weights.

```yaml
load_balancer:
  strategy: "weighted-round-robin"
backends:
  - url: "http://powerful-server:8080"
    weight: 3
  - url: "http://normal-server:8080"
    weight: 1
```

### Smooth Weighted Round Robin
Advanced weighted distribution with smooth traffic flow.

```yaml
load_balancer:
  strategy: "smooth-weighted-round-robin"
```

### Least Connections
Routes to the backend with fewest active connections.

```yaml
load_balancer:
  strategy: "least-connections"
```

### Random
Randomly selects a backend for each request.

```yaml
load_balancer:
  strategy: "random"
```

## 🏥 Health Checking

### Configuration Options

```yaml
health_check:
  enabled: true
  interval: 10s        # How often to check
  timeout: 5s          # Request timeout
  path: "/health"      # Health endpoint
  method: "GET"        # HTTP method
```

### Backend Health Endpoints

Implement health endpoints in your backends:

**Express.js Example:**
```javascript
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'healthy',
    timestamp: Date.now()
  });
});
```

**Go Example:**
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
    })
})
```

### ⚠️ Important: Health Checking Best Practices

```yaml
# ✅ Recommended for production environments
health_check:
  enabled: true  # Enables automatic recovery
  interval: 30s
  timeout: 5s
  path: "/health"
```

**Why Health Checking Matters:**
- ✅ **Auto-recovery**: Failed backends automatically come back online
- ⚡ **Zero downtime**: Traffic routing adapts to backend status
- 📊 **Observability**: Real-time backend health monitoring

## 📊 Monitoring & Metrics

### Request Logging

The load balancer provides comprehensive logging:

```
🚀 Load Balancer running on :8080
📊 Strategy: smooth-weighted-round-robin
🏢 Backends: 3 configured
✅ Added backend: http://localhost:3001 (weight: 1)
✅ Added backend: http://localhost:3002 (weight: 2)
🏥 Health checking enabled (interval: 10s)
🎯 /api/users -> http://localhost:3002
✅ Backend http://localhost:3001 healthy (latency: 2ms)
```

### Health Status

Monitor backend health in real-time through logs:

```
✅ Backend http://localhost:3001 healthy (latency: 5ms)
✅ Backend http://localhost:3002 healthy (latency: 3ms)
❌ Backend http://localhost:3003 unhealthy: connection refused
```

## 🏗️ Architecture

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

### Core Components

- **Load Balancer**: Main orchestrator and traffic distributor
- **Strategy**: Pluggable algorithms for backend selection
- **Health Checker**: Monitors backend availability and recovery
- **Server Pool**: Manages backend lifecycle and state
- **Proxy Handler**: HTTP request forwarding and error handling

## 🛠️ Development

### Project Structure

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

### Running the Application

```bash
# Development mode
go run ./cmd/server/main.go

# With custom config
go run ./cmd/server/main.go -config=custom.yaml

# Build and run
go build -o load-balancer ./cmd/server
./load-balancer
```

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-w -s" -o load-balancer ./cmd/server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o load-balancer-linux ./cmd/server
```

### Roadmap

- [ ] Unit tests implementation
- [ ] Metrics endpoint (`/metrics`)
- [ ] Docker Compose examples
- [ ] Performance benchmarks
- [ ] Admin API for runtime configuration
