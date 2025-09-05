package health

import (
	"time"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
)

type Status int

const (
	StatusHealthy Status = iota
	StatusUnhealthy
	StatusUnknown
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

type Result struct {
	Backend   *backend.Backend
	Status    Status
	Latency   time.Duration
	Error     error
	CheckedAt time.Time
}

type Checker interface {
	Check(backend *backend.Backend) *Result
	StartMonitoring(backends []*backend.Backend)
	StopMonitoring()
	GetResults() map[string]*Result
}
