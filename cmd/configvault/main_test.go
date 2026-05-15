package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain_VersionFlag(t *testing.T) {
	if os.Getenv("RUN_SUBPROCESS") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_VersionFlag", "-version")
	cmd.Env = append(os.Environ(), "RUN_SUBPROCESS=1")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected clean exit, got: %v", err)
	}

	got := string(out)
	if got == "" {
		t.Error("expected version output, got empty string")
	}
}

func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("RUN_SUBPROCESS") == "1" {
		main()
		return
	}

	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "missing.yaml")

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_MissingConfig", "-config", nonExistent)
	cmd.Env = append(os.Environ(), "RUN_SUBPROCESS=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil")
	}
}

func TestVersion_Constant(t *testing.T) {
	if version == "" {
		t.Error("version constant should not be empty")
	}
}

func TestMain_ConfigFlag_Default(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	configPath := fs.String("config", "configvault.yaml", "")
	_ = fs.Parse([]string{})

	if *configPath != "configvault.yaml" {
		t.Errorf("expected default config path 'configvault.yaml', got %q", *configPath)
	}
}
