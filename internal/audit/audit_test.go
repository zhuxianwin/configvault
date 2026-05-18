package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/configvault/internal/audit"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestLogger_Log_WritesEntry(t *testing.T) {
	path := tempLog(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()

	entry := audit.Entry{
		Event:    audit.EventSync,
		Provider: "vault",
		Target:   ".env",
		Keys:     []string{"DB_URL", "API_KEY"},
		Success:  true,
	}
	if err := l.Log(entry); err != nil {
		t.Fatalf("Log: %v", err)
	}
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	var got audit.Entry
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected one log line")
	}
	if err := json.Unmarshal(scanner.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Event != audit.EventSync {
		t.Errorf("event = %q, want %q", got.Event, audit.EventSync)
	}
	if got.Provider != "vault" {
		t.Errorf("provider = %q, want vault", got.Provider)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_Log_SetsTimestampAutomatically(t *testing.T) {
	path := tempLog(t)
	l, _ := audit.New(path)
	defer l.Close()

	before := time.Now().UTC()
	_ = l.Log(audit.Entry{Event: audit.EventRead, Provider: "ssm", Target: ".env", Success: true})
	after := time.Now().UTC()
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	var got audit.Entry
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	json.Unmarshal(scanner.Bytes(), &got)

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestLogger_Log_MultipleEntries(t *testing.T) {
	path := tempLog(t)
	l, _ := audit.New(path)
	defer l.Close()

	for i := 0; i < 3; i++ {
		_ = l.Log(audit.Entry{Event: audit.EventWrite, Provider: "vault", Target: ".env", Success: true})
	}
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}
