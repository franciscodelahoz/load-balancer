package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL               *url.URL
	Alive             bool
	ReverseProxy      *httputil.ReverseProxy
	RequestsCount     uint64
	ErrorCount        uint64
	LastErrorTime     time.Time
	ResponseTimeSum   time.Duration
	Weight            uint64
	MaxConnections    uint64
	activeConnections uint64
	mutex             sync.RWMutex
}

func CreateBackendInstance(url url.URL, weight uint64, maxConnections uint64) *Backend {
	return &Backend{
		URL:            &url,
		ReverseProxy:   httputil.NewSingleHostReverseProxy(&url),
		Alive:          true,
		Weight:         weight,
		MaxConnections: maxConnections,
	}
}

func (b *Backend) SetHealth(healthy bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Alive = healthy
}

func (b *Backend) IsAlive() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.Alive
}

func (b *Backend) IncrementRequestsCount() {
	atomic.AddUint64(&b.RequestsCount, 1)
}

func (b *Backend) IncrementErrorCount() {
	atomic.AddUint64(&b.ErrorCount, 1)
	b.LastErrorTime = time.Now()
}

func (b *Backend) GetRequestsCount() uint64 {
	return atomic.LoadUint64(&b.RequestsCount)
}

func (b *Backend) GetErrorCount() uint64 {
	return atomic.LoadUint64(&b.ErrorCount)
}

func (b *Backend) IncrementActiveConnections() {
	atomic.AddUint64(&b.activeConnections, 1)
}

func (b *Backend) DecrementActiveConnections() {
	atomic.AddUint64(&b.activeConnections, ^uint64(0))
}

func (b *Backend) GetActiveConnectionsCount() uint64 {
	return atomic.LoadUint64(&b.activeConnections)
}

func (b *Backend) GetWeight() uint64 {
	return atomic.LoadUint64(&b.Weight)
}
