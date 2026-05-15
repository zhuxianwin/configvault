package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/configvault/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_VaultConfig(t *testing.T) {
	path := writeTemp(t, `
provider: vault
output: .env
vault:
  address: http://localhost:8200
  token: root
  path: secret/data/app
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != config.ProviderVault {
		t.Errorf("provider = %q; want %q", cfg.Provider, config.ProviderVault)
	}
	if cfg.Vault.Path != "secret/data/app" {
		t.Errorf("vault.path = %q; want %q", cfg.Vault.Path, "secret/data/app")
	}
}

func TestLoad_SSMConfig(t *testing.T) {
	path := writeTemp(t, `
provider: ssm
output: .env.local
ssm:
  region: us-east-1
  path: /myapp/prod
  decrypt: true
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != config.ProviderSSM {
		t.Errorf("provider = %q; want %q", cfg.Provider, config.ProviderSSM)
	}
	if !cfg.SSM.Decrypt {
		t.Error("ssm.decrypt should be true")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_UnknownProvider(t *testing.T) {
	path := writeTemp(t, `
provider: unknown
output: .env
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown provider")
	}
}

func TestLoad_MissingOutput(t *testing.T) {
	path := writeTemp(t, `
provider: vault
vault:
  address: http://localhost:8200
  token: root
  path: secret/data/app
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing output")
	}
}

func TestLoad_MissingVaultSection(t *testing.T) {
	path := writeTemp(t, `
provider: vault
output: .env
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing vault section")
	}
}
