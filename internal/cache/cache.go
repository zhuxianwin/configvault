package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached set of secrets with metadata.
type Entry struct {
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired returns true if the cache entry is older than its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return true
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache handles reading and writing cached secret entries to disk.
type Cache struct {
	dir string
}

// New creates a new Cache that stores entries under the given directory.
func New(dir string) *Cache {
	return &Cache{dir: dir}
}

// cacheFilePath returns the path to the cache file for a given key.
func (c *Cache) cacheFilePath(key string) string {
	return filepath.Join(c.dir, key+".json")
}

// Get reads a cached entry for the given key.
// Returns nil, nil if the file does not exist.
func (c *Cache) Get(key string) (*Entry, error) {
	data, err := os.ReadFile(c.cacheFilePath(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Set writes a cache entry for the given key.
func (c *Cache) Set(key string, secrets map[string]string, ttl time.Duration) error {
	if err := os.MkdirAll(c.dir, 0o700); err != nil {
		return err
	}
	entry := Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       ttl,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.cacheFilePath(key), data, 0o600)
}

// Invalidate removes the cache entry for the given key.
func (c *Cache) Invalidate(key string) error {
	err := os.Remove(c.cacheFilePath(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
