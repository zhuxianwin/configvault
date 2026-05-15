package sync

import (
	"fmt"
	"log"

	"github.com/configvault/internal/dotenv"
	"github.com/configvault/internal/provider"
)

// Syncer fetches secrets from a provider and writes them to a dotenv file.
type Syncer struct {
	provider provider.Provider
	writer   *dotenv.Writer
	reader   *dotenv.Reader
	outPath  string
}

// New creates a new Syncer.
func New(p provider.Provider, outPath string) *Syncer {
	return &Syncer{
		provider: p,
		writer:   dotenv.NewWriter(outPath),
		reader:   dotenv.NewReader(outPath),
		outPath:  outPath,
	}
}

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Added   int
	Updated int
	Total   int
}

// Sync fetches secrets from the provider and merges them into the dotenv file.
// Existing keys not present in the provider are preserved.
func (s *Syncer) Sync() (*SyncResult, error) {
	existing, err := s.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading existing dotenv: %w", err)
	}

	fetched, err := s.provider.GetSecrets()
	if err != nil {
		return nil, fmt.Errorf("fetching secrets from %s: %w", s.provider.Name(), err)
	}

	result := &SyncResult{}
	merged := make(map[string]string, len(existing))

	for k, v := range existing {
		merged[k] = v
	}

	for k, v := range fetched {
		if _, exists := existing[k]; !exists {
			result.Added++
		} else if existing[k] != v {
			result.Updated++
		}
		merged[k] = v
	}

	result.Total = len(merged)

	if err := s.writer.Write(merged); err != nil {
		return nil, fmt.Errorf("writing dotenv file: %w", err)
	}

	log.Printf("[%s] sync complete: %d added, %d updated, %d total",
		s.provider.Name(), result.Added, result.Updated, result.Total)

	return result, nil
}
