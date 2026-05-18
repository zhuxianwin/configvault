package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ProviderConfig holds provider-specific settings.
type ProviderConfig struct {
	Type string `yaml:"type"` // "vault" or "ssm"

	// Vault options
	Address string `yaml:"address,omitempty"`
	Token   string `yaml:"token,omitempty"`
	Path    string `yaml:"path,omitempty"`

	// SSM options
	Region    string `yaml:"region,omitempty"`
	SSMPath   string `yaml:"ssm_path,omitempty"`
	Decrypt   bool   `yaml:"decrypt,omitempty"`
}

// CacheConfig controls local caching behaviour.
type CacheConfig struct {
	Enabled bool          `yaml:"enabled"`
	Dir     string        `yaml:"dir"`
	TTL     time.Duration `yaml:"ttl"`
}

// Config is the top-level configvault configuration.
type Config struct {
	OutputFile string         `yaml:"output_file"`
	Provider   ProviderConfig `yaml:"provider"`
	Cache      CacheConfig    `yaml:"cache"`
}

// Load reads and validates a Config from the YAML file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}
	if cfg.Cache.Dir == "" {
		cfg.Cache.Dir = ".configvault-cache"
	}
	if cfg.Cache.TTL == 0 {
		cfg.Cache.TTL = 5 * time.Minute
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	switch cfg.Provider.Type {
	case "vault":
		if cfg.Provider.Address == "" {
			return fmt.Errorf("config: vault provider requires 'address'")
		}
		if cfg.Provider.Path == "" {
			return fmt.Errorf("config: vault provider requires 'path'")
		}
	case "ssm":
		if cfg.Provider.SSMPath == "" {
			return fmt.Errorf("config: ssm provider requires 'ssm_path'")
		}
	case "":
		return fmt.Errorf("config: 'provider.type' is required")
	default:
		return fmt.Errorf("config: unknown provider type %q", cfg.Provider.Type)
	}
	return nil
}
