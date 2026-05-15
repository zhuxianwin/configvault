package sync

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/configvault/internal/diff"
	"github.com/configvault/internal/dotenv"
	"github.com/configvault/internal/provider"
)

// Syncer fetches secrets from a provider and writes them to a dotenv file.
type Syncer struct {
	provider provider.Provider
	writer   *dotenv.Writer
	reader   *dotenv.Reader
	output   io.Writer
}

// New creates a new Syncer.
func New(p provider.Provider, w *dotenv.Writer, r *dotenv.Reader) *Syncer {
	return &Syncer{
		provider: p,
		writer:   w,
		reader:   r,
		output:   os.Stdout,
	}
}

// Sync fetches secrets and merges them into the target dotenv file.
// It prints a diff summary of changes to the configured output.
func (s *Syncer) Sync(ctx context.Context, path string) error {
	secrets, err := s.provider.GetSecrets(ctx)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	existing, err := s.reader.Read(path)
	if err != nil {
		return fmt.Errorf("reading existing dotenv: %w", err)
	}

	changes := diff.Compute(existing, secrets)
	diff.Fprint(s.output, changes, false)

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range secrets {
		merged[k] = v
	}

	if err := s.writer.Write(path, merged); err != nil {
		return fmt.Errorf("writing dotenv: %w", err)
	}

	return nil
}
