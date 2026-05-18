package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventSync   EventType = "sync"
	EventRead   EventType = "read"
	EventWrite  EventType = "write"
	EventRotate EventType = "rotate"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Provider  string    `json:"provider"`
	Target    string    `json:"target"`
	Keys      []string  `json:"keys,omitempty"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit entries to a file.
type Logger struct {
	path string
	f    *os.File
}

// New creates a new audit Logger that appends to the given file path.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{path: path, f: f}, nil
}

// Log writes an Entry to the audit log as a JSON line.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.f, "%s\n", data)
	return err
}

// Close closes the underlying log file.
func (l *Logger) Close() error {
	return l.f.Close()
}
