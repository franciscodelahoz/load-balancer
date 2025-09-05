package loadbalancer

import (
	"net/http"
	"sync"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/health"
)

type LoadBalancer struct {
	serverPool *ServerPool
	strategy   LoadBalancerStrategy
	health     *health.HealthChecker
	mutex      sync.RWMutex
}

func NewLoadBalancer(strategy LoadBalancerStrategy) *LoadBalancer {
	return &LoadBalancer{
		serverPool: NewServerPool(),
		strategy:   strategy,
	}
}

func (lb *LoadBalancer) AddBackend(backend *backend.Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.serverPool.AddBackend(backend)

	if lb.health != nil {
		lb.health.RegisterBackend(backend)
	}

	if eventAware, ok := lb.strategy.(BackendEventAware); ok {
		eventAware.OnBackendAdded(backend)
	}
}

func (lb *LoadBalancer) StartHealthChecking(config health.Config) {
	lb.health = health.NewHealthChecker(&config)

	lb.mutex.RLock()
	backends := lb.serverPool.GetAllBackends()
	lb.mutex.RUnlock()

	for _, backend := range backends {
		lb.health.RegisterBackend(backend)
	}

	lb.health.Start()
}

func (lb *LoadBalancer) StopHealthChecking() {
	if lb.health != nil {
		lb.health.Stop()
	}
}

func (lb *LoadBalancer) GetNextBackend(r *http.Request) *backend.Backend {
	selectedBackend := lb.strategy.GetNextBackend(lb.serverPool, r)

	if selectedBackend != nil {
		selectedBackend.IncrementActiveConnections()
	}

	return selectedBackend
}

func (lb *LoadBalancer) OnRequestCompleted(backend *backend.Backend) {
	if backend != nil {
		backend.DecrementActiveConnections()
	}
}

func (lb *LoadBalancer) GetStrategyName() string {
	return lb.strategy.GetStrategyName()
}
