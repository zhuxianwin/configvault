package sync_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/example/configvault/internal/audit"
	"github.com/example/configvault/internal/dotenv"
	csync "github.com/example/configvault/internal/sync"
)

type mockProvider struct {
	name    string
	secrets map[string]string
	err     error
}

func (m *mockProvider) GetSecrets(_ context.Context) (map[string]string, error) {
	return m.secrets, m.err
}
func (m *mockProvider) Name() string { return m.name }

func TestSyncer_Sync_WritesSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	p := &mockProvider{name: "vault", secrets: map[string]string{"FOO": "bar", "BAZ": "qux"}}
	s := csync.New(p, dotenv.NewReader(), dotenv.NewWriter(), &audit.NoopLogger{})

	d, err := s.Sync(context.Background(), path)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if len(d.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(d.Added))
	}

	reader := dotenv.NewReader()
	got, err := reader.Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("FOO = %q, want bar", got["FOO"])
	}
}

func TestSyncer_Sync_MergesWithExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := dotenv.NewWriter()
	_ = w.Write(path, map[string]string{"EXISTING": "value", "FOO": "old"})

	p := &mockProvider{name: "ssm", secrets: map[string]string{"FOO": "new"}}
	s := csync.New(p, dotenv.NewReader(), w, nil)

	_, err := s.Sync(context.Background(), path)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	got, _ := dotenv.NewReader().Read(path)
	if got["EXISTING"] != "value" {
		t.Errorf("EXISTING = %q, want value", got["EXISTING"])
	}
	if got["FOO"] != "new" {
		t.Errorf("FOO = %q, want new", got["FOO"])
	}
}

func TestSyncer_Sync_ProviderError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	p := &mockProvider{name: "vault", err: errors.New("connection refused")}
	s := csync.New(p, dotenv.NewReader(), dotenv.NewWriter(), nil)

	_, err := s.Sync(context.Background(), path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
