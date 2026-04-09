package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// PortBinding represents a port that is currently bound
type PortBinding struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"` // tcp or udp
	Process   string    `json:"process"`
	PID       int       `json:"pid"`
	Timestamp time.Time `json:"timestamp"`
}

// Scanner handles port scanning operations
type Scanner struct {
	minPort int
	maxPort int
}

// New creates a new Scanner with the specified port range
func New(minPort, maxPort int) *Scanner {
	if minPort < 1 {
		minPort = 1
	}
	if maxPort > 65535 {
		maxPort = 65535
	}
	return &Scanner{
		minPort: minPort,
		maxPort: maxPort,
	}
}

// ScanPorts scans for active port bindings in the configured range
func (s *Scanner) ScanPorts() ([]PortBinding, error) {
	var bindings []PortBinding

	// Scan TCP ports
	for port := s.minPort; port <= s.maxPort; port++ {
		if s.isPortOpen(port, "tcp") {
			binding := PortBinding{
				Port:      port,
				Protocol:  "tcp",
				Timestamp: time.Now(),
			}
			bindings = append(bindings, binding)
		}
	}

	// Scan UDP ports (common ones)
	for port := s.minPort; port <= s.maxPort; port++ {
		if s.isPortOpen(port, "udp") {
			binding := PortBinding{
				Port:      port,
				Protocol:  "udp",
				Timestamp: time.Now(),
			}
			bindings = append(bindings, binding)
		}
	}

	return bindings, nil
}

// isPortOpen checks if a specific port is open
func (s *Scanner) isPortOpen(port int, protocol string) bool {
	address := net.JoinHostPort("localhost", strconv.Itoa(port))
	conn, err := net.DialTimeout(protocol, address, 100*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// GetBindingKey returns a unique key for a port binding
func (pb *PortBinding) GetBindingKey() string {
	return fmt.Sprintf("%s:%d", strings.ToLower(pb.Protocol), pb.Port)
}
