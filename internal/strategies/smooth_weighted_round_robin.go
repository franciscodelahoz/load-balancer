package strategies

import (
	"net/http"
	"sync"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type BackendWeight struct {
	Weight        uint64
	CurrentWeight int64
}

type SmoothWeightedRoundRobin struct {
	backendWeights map[*backend.Backend]*BackendWeight
	mutex          sync.RWMutex
}

func NewSmoothWeightedRoundRobin() *SmoothWeightedRoundRobin {
	return &SmoothWeightedRoundRobin{
		backendWeights: make(map[*backend.Backend]*BackendWeight),
	}
}

func (swrr *SmoothWeightedRoundRobin) OnBackendAdded(backend *backend.Backend) {
	swrr.mutex.Lock()
	defer swrr.mutex.Unlock()

	swrr.backendWeights[backend] = &BackendWeight{
		Weight:        backend.GetWeight(),
		CurrentWeight: 0,
	}
}

func (swrr *SmoothWeightedRoundRobin) OnBackendRemoved(backend *backend.Backend) {
	swrr.mutex.Lock()
	defer swrr.mutex.Unlock()

	delete(swrr.backendWeights, backend)
}

func (swrr *SmoothWeightedRoundRobin) OnBackendWeightChanged(backend *backend.Backend, oldWeight, newWeight uint64) {
	swrr.mutex.Lock()
	defer swrr.mutex.Unlock()

	if bw, exists := swrr.backendWeights[backend]; exists {
		bw.Weight = newWeight
	}
}

func (swrr *SmoothWeightedRoundRobin) GetNextBackend(pool *loadbalancer.ServerPool, r *http.Request) *backend.Backend {
	var aliveBackends []*backend.Backend = pool.GetAliveBackends()

	if len(aliveBackends) == 0 {
		return nil
	}

	swrr.mutex.Lock()
	defer swrr.mutex.Unlock()

	var selectedBackend *backend.Backend
	var maxCurrentWeight int64 = -1
	var totalWeight int64 = 0

	for _, backend := range aliveBackends {
		backendWeight, exists := swrr.backendWeights[backend]

		if !exists {
			continue
		}

		backendWeight.CurrentWeight += int64(backendWeight.Weight)
		totalWeight += int64(backendWeight.Weight)

		if backendWeight.CurrentWeight > maxCurrentWeight {
			maxCurrentWeight = backendWeight.CurrentWeight
			selectedBackend = backend
		}
	}

	if selectedBackend != nil {
		swrr.backendWeights[selectedBackend].CurrentWeight -= totalWeight
	}

	return selectedBackend
}

func (swrr *SmoothWeightedRoundRobin) GetStrategyName() string {
	return "Smooth Weighted Round Robin"
}
