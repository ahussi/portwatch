package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ProcessInfo contains information about a process using a port
type ProcessInfo struct {
	PID  int
	Name string
}

// GetProcessByPort attempts to find the process using a specific port
func GetProcessByPort(port int, protocol string) (*ProcessInfo, error) {
	switch runtime.GOOS {
	case "linux":
		return getProcessLinux(port, protocol)
	case "darwin":
		return getProcessDarwin(port, protocol)
	case "windows":
		return getProcessWindows(port, protocol)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getProcessLinux finds process info on Linux using lsof or ss
func getProcessLinux(port int, protocol string) (*ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf("%s:%d", protocol, port), "-t")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute lsof: %w", err)
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return nil, fmt.Errorf("no process found")
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PID: %w", err)
	}

	// Get process name
	cmdName := exec.Command("ps", "-p", pidStr, "-o", "comm=")
	nameOutput, _ := cmdName.Output()
	name := strings.TrimSpace(string(nameOutput))

	return &ProcessInfo{PID: pid, Name: name}, nil
}

// getProcessDarwin finds process info on macOS
func getProcessDarwin(port int, protocol string) (*ProcessInfo, error) {
	return getProcessLinux(port, protocol) // macOS uses similar commands
}

// getProcessWindows finds process info on Windows using netstat
func getProcessWindows(port int, protocol string) (*ProcessInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute netstat: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	portStr := fmt.Sprintf(":%d", port)

	for _, line := range lines {
		if strings.Contains(line, portStr) && strings.Contains(strings.ToUpper(line), strings.ToUpper(protocol)) {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pid, err := strconv.Atoi(fields[len(fields)-1])
				if err == nil {
					return &ProcessInfo{PID: pid, Name: "unknown"}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no process found")
}
