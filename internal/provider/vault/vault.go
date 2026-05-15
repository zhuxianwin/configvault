package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// httpClient is the subset of http.Client used for requests.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Provider retrieves secrets from HashiCorp Vault (KV v2).
type Provider struct {
	address string
	token   string
	client  httpClient
}

// New creates a new Vault Provider.
func New(address, token string, client httpClient) *Provider {
	if client == nil {
		client = http.DefaultClient
	}
	return &Provider{address: strings.TrimRight(address, "/"), token: token, client: client}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "vault"
}

// GetSecrets fetches secrets from a KV v2 path (e.g. "secret/data/myapp").
func (p *Provider) GetSecrets(ctx context.Context, path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", p.address, strings.TrimPrefix(path, "/"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vault: unexpected status %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("vault: decode response: %w", err)
	}

	return result.Data.Data, nil
}
