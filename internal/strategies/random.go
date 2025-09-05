package strategies

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type RandomStrategy struct{}

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{}
}

func (rs *RandomStrategy) GetNextBackend(pool *loadbalancer.ServerPool, r *http.Request) *backend.Backend {
	var aliveBackends []*backend.Backend = pool.GetAliveBackends()

	if len(aliveBackends) == 0 {
		return nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(aliveBackends))

	return aliveBackends[index]
}

func (rs *RandomStrategy) GetStrategyName() string {
	return "Random"
}
