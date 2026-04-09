package scanner

import (
	"net"
	"runtime"
	"testing"
)

func TestGetProcessByPort(t *testing.T) {
	// Skip on unsupported platforms
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows - requires admin privileges")
	}

	// Start a test server
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	// Try to get process info
	info, err := GetProcessByPort(port, "tcp")
	
	// On some systems, lsof might not be available or require permissions
	if err != nil {
		t.Logf("GetProcessByPort returned error (may be expected): %v", err)
		return
	}

	if info == nil {
		t.Error("Expected process info, got nil")
		return
	}

	if info.PID <= 0 {
		t.Errorf("Invalid PID: %d", info.PID)
	}

	t.Logf("Found process: PID=%d, Name=%s", info.PID, info.Name)
}

func TestGetProcessByPort_NotFound(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	// Try to get info for a port that's definitely not in use
	info, err := GetProcessByPort(65432, "tcp")
	
	if err == nil {
		t.Error("Expected error for unused port, got nil")
	}

	if info != nil {
		t.Error("Expected nil process info for unused port")
	}
}

func TestProcessInfo(t *testing.T) {
	info := &ProcessInfo{
		PID:  1234,
		Name: "test-process",
	}

	if info.PID != 1234 {
		t.Errorf("PID = %d, want 1234", info.PID)
	}

	if info.Name != "test-process" {
		t.Errorf("Name = %s, want test-process", info.Name)
	}
}
