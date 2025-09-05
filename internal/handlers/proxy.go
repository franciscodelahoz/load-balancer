package handlers

import (
	"log"
	"net/http"

	"github.com/franciscodelahoz/load-balancer/internal/backend"
	"github.com/franciscodelahoz/load-balancer/internal/loadbalancer"
)

type ProxyHandler struct {
	loadBalancer *loadbalancer.LoadBalancer
}

func NewProxyHandler(lb *loadbalancer.LoadBalancer) *ProxyHandler {
	return &ProxyHandler{
		loadBalancer: lb,
	}
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}

	return "http"
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	selectedBackend := ph.loadBalancer.GetNextBackend(r)

	if selectedBackend == nil {
		log.Printf("❌ No healthy backends available")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	selectedBackend.IncrementRequestsCount()
	selectedBackend.IncrementActiveConnections()

	defer func() {
		selectedBackend.DecrementActiveConnections()
		ph.loadBalancer.OnRequestCompleted(selectedBackend)
	}()

	log.Printf("🎯 %s -> %s", r.URL.Path, selectedBackend.URL.String())

	ph.proxyRequest(w, r, selectedBackend)
}

func (ph *ProxyHandler) proxyRequest(w http.ResponseWriter, r *http.Request, b *backend.Backend) {
	proxy := b.ReverseProxy

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		originalHost := req.Host
		req.Host = b.URL.Host

		req.Header.Set("X-Forwarded-Host", originalHost)
		req.Header.Set("X-Forwarded-Proto", getScheme(r))
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("❌ Proxy error for backend %s: %v", b.URL.String(), err)

		b.IncrementErrorCount()

		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode >= 400 {
			b.IncrementErrorCount()
		}
		return nil
	}

	proxy.ServeHTTP(w, r)
}
