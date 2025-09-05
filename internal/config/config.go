package config

import (
	"fmt"
	"os"

	"github.com/franciscodelahoz/load-balancer/internal/health"
	"gopkg.in/yaml.v3"
)

func (cfg *Config) applyDefaults() {
	// Server defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = DefaultPort
	}

	// LoadBalancer defaults
	if cfg.LoadBalancer.Strategy == "" {
		cfg.LoadBalancer.Strategy = DefaultStrategy
	}

	// HealthCheck defaults
	if cfg.HealthCheck.Enabled == nil {
		enabled := DefaultEnabled
		cfg.HealthCheck.Enabled = &enabled
	}
	if cfg.HealthCheck.Interval == 0 {
		cfg.HealthCheck.Interval = DefaultInterval
	}
	if cfg.HealthCheck.Timeout == 0 {
		cfg.HealthCheck.Timeout = DefaultTimeout
	}
	if cfg.HealthCheck.Path == "" {
		cfg.HealthCheck.Path = DefaultPath
	}
	if cfg.HealthCheck.Method == "" {
		cfg.HealthCheck.Method = DefaultMethod
	}

	// Backend defaults
	for i := range cfg.Backends {
		if cfg.Backends[i].Weight == 0 {
			cfg.Backends[i].Weight = DefaultWeight
		}
	}
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.applyDefaults()
	return config, nil
}

func (cfg *Config) GetHealthConfig() *health.Config {
	return &health.Config{
		Interval: cfg.HealthCheck.Interval,
		Timeout:  cfg.HealthCheck.Timeout,
		Path:     cfg.HealthCheck.Path,
		Method:   cfg.HealthCheck.Method,
	}
}

func (cfg *Config) IsHealthCheckEnabled() bool {
	if cfg.HealthCheck.Enabled == nil {
		return DefaultEnabled
	}
	return *cfg.HealthCheck.Enabled
}
