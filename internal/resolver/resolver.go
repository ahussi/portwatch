// Package resolver maps port numbers to well-known service names.
package resolver

import (
	"fmt"
	"sync"
)

// ServiceInfo holds metadata about a well-known service.
type ServiceInfo struct {
	Name     string
	Protocol string
	Port     int
}

// Resolver resolves port numbers to service names.
type Resolver struct {
	mu       sync.RWMutex
	services map[int]*ServiceInfo
}

// New returns a Resolver pre-loaded with common well-known ports.
func New() *Resolver {
	r := &Resolver{
		services: make(map[int]*ServiceInfo),
	}
	r.loadDefaults()
	return r
}

// Resolve returns the ServiceInfo for the given port, or nil if unknown.
func (r *Resolver) Resolve(port int) *ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.services[port]
}

// Register adds or overwrites a service entry for the given port.
func (r *Resolver) Register(port int, name, protocol string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[port] = &ServiceInfo{Name: name, Protocol: protocol, Port: port}
}

// Name returns the service name for port, or a formatted fallback.
func (r *Resolver) Name(port int) string {
	if svc := r.Resolve(port); svc != nil {
		return svc.Name
	}
	return fmt.Sprintf("port/%d", port)
}

func (r *Resolver) loadDefaults() {
	defaults := []ServiceInfo{
		{Port: 21, Name: "ftp", Protocol: "tcp"},
		{Port: 22, Name: "ssh", Protocol: "tcp"},
		{Port: 23, Name: "telnet", Protocol: "tcp"},
		{Port: 25, Name: "smtp", Protocol: "tcp"},
		{Port: 53, Name: "dns", Protocol: "udp"},
		{Port: 80, Name: "http", Protocol: "tcp"},
		{Port: 110, Name: "pop3", Protocol: "tcp"},
		{Port: 143, Name: "imap", Protocol: "tcp"},
		{Port: 443, Name: "https", Protocol: "tcp"},
		{Port: 3306, Name: "mysql", Protocol: "tcp"},
		{Port: 5432, Name: "postgres", Protocol: "tcp"},
		{Port: 6379, Name: "redis", Protocol: "tcp"},
		{Port: 8080, Name: "http-alt", Protocol: "tcp"},
		{Port: 27017, Name: "mongodb", Protocol: "tcp"},
	}
	for i := range defaults {
		svc := defaults[i]
		r.services[svc.Port] = &svc
	}
}
