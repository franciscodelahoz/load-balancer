package strategies

import (
	"net/http"
	"sync/atomic"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type WeightedRoundRobin struct {
	currentIndex uint64
	currentCount uint64
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		currentIndex: 0,
		currentCount: 0,
	}
}

func (wrr *WeightedRoundRobin) GetNextBackend(pool *loadbalancer.ServerPool, r *http.Request) *backend.Backend {
	var aliveBackends []*backend.Backend = pool.GetAliveBackends()

	if len(aliveBackends) == 0 {
		return nil
	}

	currentIndex := atomic.LoadUint64(&wrr.currentIndex)
	currentCount := atomic.LoadUint64(&wrr.currentCount)

	normalizedIndex := currentIndex % uint64(len(aliveBackends))
	currentBackend := aliveBackends[normalizedIndex]

	weight := currentBackend.GetWeight()

	if currentCount < weight {
		atomic.AddUint64(&wrr.currentCount, 1)
		return currentBackend
	}

	nextBackendIndex := (normalizedIndex + 1) % uint64(len(aliveBackends))

	atomic.StoreUint64(&wrr.currentCount, 1)
	atomic.StoreUint64(&wrr.currentIndex, nextBackendIndex)

	return aliveBackends[nextBackendIndex]
}

func (wrr *WeightedRoundRobin) GetStrategyName() string {
	return "Weighted Round Robin"
}
