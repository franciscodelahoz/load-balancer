package loadbalancer

import (
	"net/http"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
)

type LoadBalancerStrategy interface {
	GetNextBackend(pool *ServerPool, r *http.Request) *backend.Backend
	GetStrategyName() string
}

type BackendEventAware interface {
	OnBackendAdded(backend *backend.Backend)
	OnBackendRemoved(backend *backend.Backend)
	OnBackendWeightChanged(backend *backend.Backend, oldWeight, newWeight uint64)
}
