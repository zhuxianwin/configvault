package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/configvault/internal/cache"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestCache_SetAndGet_ReturnsEntry(t *testing.T) {
	c := cache.New(tempDir(t))
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	if err := c.Set("myapp", secrets, 10*time.Minute); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	entry, err := c.Get("myapp")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", entry.Secrets["DB_HOST"])
	}
}

func TestCache_Get_NonExistent_ReturnsNil(t *testing.T) {
	c := cache.New(tempDir(t))

	entry, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry, got %+v", entry)
	}
}

func TestCache_IsExpired_TTLElapsed(t *testing.T) {
	entry := &cache.Entry{
		Secrets:   map[string]string{"KEY": "val"},
		FetchedAt: time.Now().Add(-5 * time.Minute),
		TTL:       1 * time.Minute,
	}
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}
}

func TestCache_IsExpired_WithinTTL(t *testing.T) {
	entry := &cache.Entry{
		Secrets:   map[string]string{"KEY": "val"},
		FetchedAt: time.Now(),
		TTL:       10 * time.Minute,
	}
	if entry.IsExpired() {
		t.Error("expected entry to not be expired")
	}
}

func TestCache_Invalidate_RemovesEntry(t *testing.T) {
	c := cache.New(tempDir(t))

	_ = c.Set("myapp", map[string]string{"X": "1"}, time.Minute)
	if err := c.Invalidate("myapp"); err != nil {
		t.Fatalf("Invalidate() error: %v", err)
	}

	entry, err := c.Get("myapp")
	if err != nil {
		t.Fatalf("Get() after invalidate error: %v", err)
	}
	if entry != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestCache_Invalidate_NonExistent_NoError(t *testing.T) {
	c := cache.New(tempDir(t))
	if err := c.Invalidate("ghost"); err != nil {
		t.Errorf("Invalidate() on missing key should not error: %v", err)
	}
}
