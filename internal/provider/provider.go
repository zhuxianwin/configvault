package provider

import "context"

// Provider defines the interface for secret backends.
type Provider interface {
	// GetSecrets retrieves secrets from the backend for the given path.
	GetSecrets(ctx context.Context, path string) (map[string]string, error)
	// Name returns the human-readable name of the provider.
	Name() string
}
