package strategies

import (
	"net/http"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type LeastConnectionsStrategy struct{}

func NewLeastConnectionsStrategy() *LeastConnectionsStrategy {
	return &LeastConnectionsStrategy{}
}

func (lc *LeastConnectionsStrategy) GetNextBackend(pool *loadbalancer.ServerPool, r *http.Request) *backend.Backend {
	var aliveBackends []*backend.Backend = pool.GetAliveBackends()

	if len(aliveBackends) == 0 {
		return nil
	}

	selectedBackend := aliveBackends[0]
	selectedBackendConnections := selectedBackend.GetActiveConnectionsCount()

	for i := 1; i < len(aliveBackends); i += 1 {
		connections := aliveBackends[i].GetActiveConnectionsCount()

		if connections < selectedBackendConnections {
			selectedBackend = aliveBackends[i]
			selectedBackendConnections = connections
		}
	}

	return selectedBackend
}

func (lc *LeastConnectionsStrategy) GetStrategyName() string {
	return "Least Connections"
}
