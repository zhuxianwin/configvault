package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_Write_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWriter_Write_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(filepath.Join(dir, ".env"))

	secrets := map[string]string{
		"SECRET_KEY": "hello world",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, ".env"))
	if !strings.Contains(string(data), `"hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", string(data))
	}
}

func TestWriter_Write_EmptySecretsNoFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	w := NewWriter(path)

	if err := w.Write(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected no file to be created for empty secrets")
	}
}

func TestWriter_Write_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(filepath.Join(dir, ".env"))

	secrets := map[string]string{
		"Z_VAR": "z",
		"A_VAR": "a",
		"M_VAR": "m",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, ".env"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if !strings.HasPrefix(lines[1], "A_VAR") {
		t.Errorf("expected A_VAR first, got %s", lines[1])
	}
	if !strings.HasPrefix(lines[len(lines)-1], "Z_VAR") {
		t.Errorf("expected Z_VAR last, got %s", lines[len(lines)-1])
	}
}
