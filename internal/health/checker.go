package health

import (
	"fmt"
	"log"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
)

type HealthChecker struct {
	config      *Config
	client      *http.Client
	backends    []*backend.Backend
	results     map[string]*Result
	mutex       sync.RWMutex
	stopChannel chan struct{}
	wg          sync.WaitGroup
	running     bool
}

func NewHealthChecker(config *Config) *HealthChecker {
	return &HealthChecker{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		backends:    make([]*backend.Backend, 0),
		results:     make(map[string]*Result),
		stopChannel: make(chan struct{}),
	}
}

func (hc *HealthChecker) RegisterBackend(backend *backend.Backend) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.backends = append(hc.backends, backend)
	log.Printf("✅ Registered backend for health checking: %s", backend.URL.String())
}

func (hc *HealthChecker) Start() {
	hc.mutex.Lock()

	if hc.running {
		hc.mutex.Unlock()
		return
	}

	hc.running = true
	hc.mutex.Unlock()

	hc.wg.Add(1)
	go hc.monitoringLoop()

	log.Printf("🏥 Health checker started with %d backends", len(hc.backends))
}

func (hc *HealthChecker) Stop() {
	hc.mutex.Lock()

	if !hc.running {
		hc.mutex.Unlock()
		return
	}

	hc.running = false
	hc.mutex.Unlock()

	close(hc.stopChannel)

	hc.wg.Wait()
	log.Println("🏥 Health checker stopped")
}

func (hc *HealthChecker) Check(backend *backend.Backend) *Result {
	start := time.Now()

	healthURL := backend.URL.String() + hc.config.Path

	req, err := http.NewRequest(hc.config.Method, healthURL, nil)

	if err != nil {
		return &Result{
			Backend:   backend,
			Status:    StatusUnhealthy,
			Error:     err,
			CheckedAt: time.Now(),
		}
	}

	resp, err := hc.client.Do(req)

	if err != nil {
		return &Result{
			Backend:   backend,
			Status:    StatusUnhealthy,
			Latency:   time.Since(start),
			Error:     err,
			CheckedAt: time.Now(),
		}
	}

	defer resp.Body.Close()

	status := StatusUnhealthy
	var statusError error = nil

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status = StatusHealthy
	} else {
		statusError = fmt.Errorf("unexpected HTTP status from backend: %d", resp.StatusCode)
	}

	return &Result{
		Backend:   backend,
		Status:    status,
		Latency:   time.Since(start),
		CheckedAt: time.Now(),
		Error:     statusError,
	}
}

func (hc *HealthChecker) StartMonitoring(backends []*backend.Backend) {
	hc.mutex.Lock()
	hc.backends = backends
	hc.mutex.Unlock()

	hc.wg.Add(1)
	go hc.monitoringLoop()

	log.Printf("✅ Health checker started with %d backends", len(backends))
}

func (hc *HealthChecker) StopMonitoring() {
	close(hc.stopChannel)
	hc.wg.Wait()
	log.Println("✅ Health checker stopped")
}

func (hc *HealthChecker) performHealthChecks() {
	hc.mutex.RLock()

	backends := make([]*backend.Backend, len(hc.backends))
	copy(backends, hc.backends)

	hc.mutex.RUnlock()

	if len(backends) == 0 {
		return
	}

	var wg sync.WaitGroup

	for _, backend := range backends {
		wg.Add(1)
		b := backend

		go func() {
			defer wg.Done()

			result := hc.Check(b)

			hc.mutex.Lock()
			hc.results[b.URL.String()] = result
			hc.mutex.Unlock()

			if result.Status == StatusHealthy {
				b.IncreaseConsecutiveSuccesses()
				b.ResetConsecutiveErrors()

				log.Printf("✅ Backend %s health check passed (latency: %s)", b.URL.String(), result.Latency)

				if b.GetConsecutiveSuccesses() >= hc.config.SuccessThreshold && !b.IsAlive() {
					b.SetHealth(true)
					log.Printf("✅ Backend %s marked as healthy after %d consecutive successes", b.URL.String(), b.GetConsecutiveSuccesses())
				}

			} else {
				b.ResetConsecutiveSuccesses()
				b.IncreaseConsecutiveErrors()

				log.Printf("❌ Backend %s health check failed: %v", b.URL.String(), result.Error)

				if b.GetConsecutiveErrors() >= hc.config.FailureThreshold && b.IsAlive() {
					b.SetHealth(false)
					log.Printf("❌ Backend %s marked as unhealthy after %d consecutive errors", b.URL.String(), b.GetConsecutiveErrors())
				}
			}
		}()
	}

	wg.Wait()
}

func (hc *HealthChecker) monitoringLoop() {
	defer hc.wg.Done()

	ticker := time.NewTicker(hc.config.Interval)
	defer ticker.Stop()

	// Initial check
	hc.performHealthChecks()

	for {
		select {
		case <-ticker.C:
			hc.performHealthChecks()
		case <-hc.stopChannel:
			return
		}
	}
}

func (hc *HealthChecker) GetResults() map[string]*Result {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	results := make(map[string]*Result)
	maps.Copy(results, hc.results)
	return results
}

func (hc *HealthChecker) UpdateBackends(backends []*backend.Backend) {
	hc.mutex.Lock()
	hc.backends = backends
	hc.mutex.Unlock()

	log.Printf("📝 Health checker updated with %d backends", len(backends))
}
