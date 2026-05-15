package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ProviderType represents the type of secrets provider.
type ProviderType string

const (
	ProviderVault ProviderType = "vault"
	ProviderSSM   ProviderType = "ssm"
)

// Config holds the top-level configvault configuration.
type Config struct {
	Provider ProviderType `yaml:"provider"`
	Output   string       `yaml:"output"`
	Vault    *VaultConfig `yaml:"vault,omitempty"`
	SSM      *SSMConfig   `yaml:"ssm,omitempty"`
}

// VaultConfig holds Vault-specific configuration.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Path    string `yaml:"path"`
}

// SSMConfig holds AWS SSM-specific configuration.
type SSMConfig struct {
	Region  string `yaml:"region"`
	Path    string `yaml:"path"`
	Decrypt bool   `yaml:"decrypt"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Provider != ProviderVault && c.Provider != ProviderSSM {
		return fmt.Errorf("config: unknown provider %q (must be \"vault\" or \"ssm\")", c.Provider)
	}
	if c.Output == "" {
		return fmt.Errorf("config: output path must not be empty")
	}
	if c.Provider == ProviderVault && c.Vault == nil {
		return fmt.Errorf("config: vault section required when provider is \"vault\"")
	}
	if c.Provider == ProviderSSM && c.SSM == nil {
		return fmt.Errorf("config: ssm section required when provider is \"ssm\"")
	}
	return nil
}
