package strategies

import (
	"fmt"

	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type StrategyFactory struct{}

func NewStrategyFactory() *StrategyFactory {
	return &StrategyFactory{}
}

func (s *StrategyFactory) CreateLoadbalancerStrategy(strategyName string) (loadbalancer.LoadBalancerStrategy, error) {
	switch strategyName {
	case "round-robin":
		return NewRoundRobinStrategy(), nil
	case "least-connections":
		return NewLeastConnectionsStrategy(), nil
	case "random":
		return NewRandomStrategy(), nil
	case "weighted-round-robin":
		return NewWeightedRoundRobin(), nil
	case "smooth-weighted-round-robin":
		return NewSmoothWeightedRoundRobin(), nil
	default:
		err := fmt.Errorf("unknown strategy type: %s", strategyName)
		return nil, err
	}
}
