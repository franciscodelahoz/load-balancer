package loadbalancer

import (
	"sync"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
)

type ServerPool struct {
	backends []*backend.Backend
	mutex    sync.RWMutex
}

func NewServerPool() *ServerPool {
	return &ServerPool{
		backends: make([]*backend.Backend, 0),
	}
}

func (pool *ServerPool) AddBackend(b *backend.Backend) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	pool.backends = append(pool.backends, b)
}

func (pool *ServerPool) RemoveBackend(backendURL string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for i, backend := range pool.backends {
		if backend.URL.String() == backendURL {
			pool.backends = append(pool.backends[:i], pool.backends[i+1:]...)
			break
		}
	}
}

func (pool *ServerPool) GetAllBackends() []*backend.Backend {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	result := make([]*backend.Backend, len(pool.backends))
	copy(result, pool.backends)

	return result
}

func (pool *ServerPool) GetAliveBackends() []*backend.Backend {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	var aliveBackends []*backend.Backend

	for _, backend := range pool.backends {
		if backend.IsAlive() {
			aliveBackends = append(aliveBackends, backend)
		}
	}

	return aliveBackends
}
