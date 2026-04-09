package alert

import (
	"fmt"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert represents a port-related alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Message   string
	Process   string
}

// Handler defines the interface for alert sinks.
type Handler interface {
	Send(a Alert) error
}

// Manager dispatches alerts to one or more handlers.
type Manager struct {
	handlers []Handler
}

// NewManager creates a new Manager with the given handlers.
func NewManager(handlers ...Handler) *Manager {
	return &Manager{handlers: handlers}
}

// Dispatch sends an alert to all registered handlers.
func (m *Manager) Dispatch(level Level, port int, process, message string) {
	a := Alert{
		Timestamp: time.Now(),
		Level:     level,
		Port:      port,
		Process:   process,
		Message:   message,
	}
	for _, h := range m.handlers {
		if err := h.Send(a); err != nil {
			fmt.Printf("alert handler error: %v\n", err)
		}
	}
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s | port=%d process=%q msg=%s",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Process,
		a.Message,
	)
}
