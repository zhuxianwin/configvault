package sync_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/configvault/internal/sync"
)

// mockProvider implements provider.Provider for testing.
type mockProvider struct {
	name    string
	secrets map[string]string
	err     error
}

func (m *mockProvider) GetSecrets() (map[string]string, error) {
	return m.secrets, m.err
}

func (m *mockProvider) Name() string {
	return m.name
}

func TestSyncer_Sync_WritesSecrets(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	p := &mockProvider{
		name:    "mock",
		secrets: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"},
	}

	s := sync.New(p, outPath)
	result, err := s.Sync()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Added != 2 {
		t.Errorf("expected 2 added, got %d", result.Added)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected dotenv file to be created")
	}
}

func TestSyncer_Sync_MergesWithExisting(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	// Write a pre-existing dotenv file.
	if err := os.WriteFile(outPath, []byte("EXISTING_KEY=existing_value\nDB_HOST=old_host\n"), 0600); err != nil {
		t.Fatal(err)
	}

	p := &mockProvider{
		name:    "mock",
		secrets: map[string]string{"DB_HOST": "new_host", "API_KEY": "secret"},
	}

	s := sync.New(p, outPath)
	result, err := s.Sync()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
	if result.Updated != 1 {
		t.Errorf("expected 1 updated, got %d", result.Updated)
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
}

func TestSyncer_Sync_ProviderError(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env")

	p := &mockProvider{
		name: "mock",
		err:  errors.New("connection refused"),
	}

	s := sync.New(p, outPath)
	_, err := s.Sync()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
