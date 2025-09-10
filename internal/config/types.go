package config

import "time"

type ServerConfig struct {
	Port int `yaml:"port,omitempty"`
}

type BackendConfig struct {
	URL    string `yaml:"url"`
	Weight uint64 `yaml:"weight,omitempty"`
}

type HealthCheckConfig struct {
	Enabled          *bool         `yaml:"enabled,omitempty"`
	Interval         time.Duration `yaml:"interval,omitempty"`
	Timeout          time.Duration `yaml:"timeout,omitempty"`
	Path             string        `yaml:"path,omitempty"`
	Method           string        `yaml:"method,omitempty"`
	SuccessThreshold int           `yaml:"success_threshold,omitempty"`
	FailureThreshold int           `yaml:"failure_threshold,omitempty"`
}

type LoadBalancerConfig struct {
	Strategy string `yaml:"strategy,omitempty"`
}

type Config struct {
	Server       ServerConfig       `yaml:"server,omitempty"`
	LoadBalancer LoadBalancerConfig `yaml:"load_balancer,omitempty"`
	Backends     []BackendConfig    `yaml:"backends,omitempty"`
	HealthCheck  HealthCheckConfig  `yaml:"health_check,omitempty"`
}

const (
	DefaultPort             = 8080
	DefaultStrategy         = "round-robin"
	DefaultEnabled          = true
	DefaultInterval         = 10 * time.Second
	DefaultTimeout          = 5 * time.Second
	DefaultPath             = "/health"
	DefaultMethod           = "GET"
	DefaultWeight           = uint64(1)
	DefaultSuccessThreshold = 3
	DefaultFailureThreshold = 3
)
