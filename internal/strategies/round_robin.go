package strategies

import (
	"net/http"
	"sync/atomic"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type RoundRobinStrategy struct {
	currentIndex uint64
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		currentIndex: 0,
	}
}

func (rr *RoundRobinStrategy) getNextIndex(totalBackends int) uint64 {
	var nextIndex = atomic.AddUint64(&rr.currentIndex, 1)
	var normalizedIndex = (nextIndex - 1) % uint64(totalBackends)

	return normalizedIndex
}

func (rr *RoundRobinStrategy) GetNextBackend(pool *loadbalancer.ServerPool, r *http.Request) *backend.Backend {
	var aliveBackends []*backend.Backend = pool.GetAliveBackends()

	if len(aliveBackends) == 0 {
		return nil
	}

	nextIndex := rr.getNextIndex(len(aliveBackends))
	return aliveBackends[nextIndex]
}

func (rr *RoundRobinStrategy) GetStrategyName() string {
	return "Round Robin"
}
