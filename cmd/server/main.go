package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/config"
	"github.com/franciscodelahoz/load-balancer/internal/handlers"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
	"github.com/franciscodelahoz/load-balancer/internal/strategies"
)

func main() {
	log.Println("🚀 Starting Load Balancer...")

	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		log.Printf("❌ Error loading config: %v", err)
		return
	}

	strategyFactory := strategies.NewStrategyFactory()
	strategy, err := strategyFactory.CreateLoadbalancerStrategy(cfg.LoadBalancer.Strategy)

	if err != nil {
		log.Fatalf("❌ Error creating strategy '%s': %v", cfg.LoadBalancer.Strategy, err)
	}

	loadBalancer := loadbalancer.NewLoadBalancer(strategy)

	for _, backendConfig := range cfg.Backends {
		backendURL, err := url.Parse(backendConfig.URL)

		if err != nil {
			log.Printf("❌ Invalid backend URL %s: %v", backendConfig.URL, err)
			continue
		}

		backend := backend.CreateBackendInstance(*backendURL, backendConfig.Weight, 1)
		loadBalancer.AddBackend(backend)

		log.Printf("✅ Added backend: %s (weight: %d)", backendConfig.URL, backendConfig.Weight)
	}

	if cfg.IsHealthCheckEnabled() {
		healthConfig := cfg.GetHealthConfig()
		loadBalancer.StartHealthChecking(*healthConfig)

		log.Printf("🏥 Health checking enabled (interval: %v)", cfg.HealthCheck.Interval)
	}

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	proxyHandler := handlers.NewProxyHandler(loadBalancer)

	log.Printf("🚀 Load Balancer running on :%s", address)
	log.Printf("📊 Strategy: %s", loadBalancer.GetStrategyName())
	log.Printf("🏢 Admin API: http://localhost:%s/admin/health", address)

	log.Fatal(http.ListenAndServe(address, proxyHandler))
}
