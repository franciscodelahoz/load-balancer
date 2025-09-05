package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"

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

func (ph *ProxyHandler) proxyRequest(w http.ResponseWriter, r *http.Request, b *backend.Backend) {
	proxy := httputil.NewSingleHostReverseProxy(b.URL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("❌ Proxy error for backend %s: %v", b.URL.String(), err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	selectedBackend := ph.loadBalancer.GetNextBackend(r)

	if selectedBackend == nil {
		log.Printf("❌ No healthy backends available")

		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	defer ph.loadBalancer.OnRequestCompleted(selectedBackend)

	log.Printf("🎯 %s -> %s", r.URL.Path, selectedBackend.URL.String())
	ph.proxyRequest(w, r, selectedBackend)
}
