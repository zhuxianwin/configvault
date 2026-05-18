package sync

import (
	"context"
	"fmt"

	"github.com/example/configvault/internal/audit"
	"github.com/example/configvault/internal/diff"
	"github.com/example/configvault/internal/dotenv"
	"github.com/example/configvault/internal/provider"
)

// Syncer fetches secrets from a provider and writes them to a dotenv file.
type Syncer struct {
	provider provider.Provider
	reader   *dotenv.Reader
	writer   *dotenv.Writer
	auditor  audit.Auditor
}

// New creates a Syncer with the given dependencies.
func New(p provider.Provider, r *dotenv.Reader, w *dotenv.Writer, a audit.Auditor) *Syncer {
	if a == nil {
		a = &audit.NoopLogger{}
	}
	return &Syncer{provider: p, reader: r, writer: w, auditor: a}
}

// Sync fetches secrets from the provider and merges them into the target dotenv file.
// It returns a Diff describing what changed.
func (s *Syncer) Sync(ctx context.Context, path string) (diff.Diff, error) {
	secrets, err := s.provider.GetSecrets(ctx)
	if err != nil {
		_ = s.auditor.Log(audit.Entry{
			Event:    audit.EventSync,
			Provider: s.provider.Name(),
			Target:   path,
			Success:  false,
			Message:  err.Error(),
		})
		return diff.Diff{}, fmt.Errorf("sync: fetch secrets: %w", err)
	}

	existing, err := s.reader.Read(path)
	if err != nil {
		return diff.Diff{}, fmt.Errorf("sync: read existing: %w", err)
	}

	changes := diff.Compute(existing, secrets)

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range secrets {
		merged[k] = v
	}

	if err := s.writer.Write(path, merged); err != nil {
		return diff.Diff{}, fmt.Errorf("sync: write dotenv: %w", err)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	_ = s.auditor.Log(audit.Entry{
		Event:    audit.EventSync,
		Provider: s.provider.Name(),
		Target:   path,
		Keys:     keys,
		Success:  true,
	})

	return changes, nil
}
