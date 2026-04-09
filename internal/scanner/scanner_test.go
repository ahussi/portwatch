package scanner

import (
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		minPort int
		maxPort int
		wantMin int
		wantMax int
	}{
		{"valid range", 8000, 9000, 8000, 9000},
		{"min too low", -1, 100, 1, 100},
		{"max too high", 1000, 70000, 1000, 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.minPort, tt.maxPort)
			if s.minPort != tt.wantMin {
				t.Errorf("minPort = %d, want %d", s.minPort, tt.wantMin)
			}
			if s.maxPort != tt.wantMax {
				t.Errorf("maxPort = %d, want %d", s.maxPort, tt.wantMax)
			}
		})
	}
}

func TestScanPorts(t *testing.T) {
	// Start a test server on a known port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	s := New(port, port)
	bindings, err := s.ScanPorts()
	if err != nil {
		t.Fatalf("ScanPorts failed: %v", err)
	}

	found := false
	for _, b := range bindings {
		if b.Port == port && b.Protocol == "tcp" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find binding on port %d", port)
	}
}

func TestGetBindingKey(t *testing.T) {
	pb := &PortBinding{
		Port:      8080,
		Protocol:  "TCP",
		Timestamp: time.Now(),
	}

	key := pb.GetBindingKey()
	expected := "tcp:8080"

	if key != expected {
		t.Errorf("GetBindingKey() = %s, want %s", key, expected)
	}
}

func TestIsPortOpen(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	s := New(1, 65535)

	if !s.isPortOpen(port, "tcp") {
		t.Errorf("Expected port %d to be open", port)
	}

	if s.isPortOpen(65534, "tcp") {
		t.Error("Expected port 65534 to be closed")
	}
}
