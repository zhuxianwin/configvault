package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/configvault/internal/config"
	"github.com/example/configvault/internal/dotenv"
	"github.com/example/configvault/internal/provider/ssm"
	"github.com/example/configvault/internal/provider/vault"
	"github.com/example/configvault/internal/sync"
)

const version = "0.1.0"

func main() {
	configPath := flag.String("config", "configvault.yaml", "path to config file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("configvault v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var provider interface {
		GetSecrets() (map[string]string, error)
		Name() string
	}

	switch cfg.Provider {
	case "vault":
		provider = vault.New(cfg.Vault.Address, cfg.Vault.Token, cfg.Vault.Path)
	case "ssm":
		provider = ssm.New(cfg.SSM.Region, cfg.SSM.Path, cfg.SSM.Decrypt)
	default:
		log.Fatalf("unknown provider: %s", cfg.Provider)
	}

	reader := dotenv.NewReader(cfg.OutputFile)
	writer := dotenv.NewWriter(cfg.OutputFile)
	syncer := sync.New(provider, reader, writer)

	if err := syncer.Sync(); err != nil {
		log.Fatalf("sync failed: %v", err)
	}

	fmt.Printf("synced secrets from %s to %s\n", provider.Name(), cfg.OutputFile)
}
