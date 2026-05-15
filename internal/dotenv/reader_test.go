package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReader_Read_ParsesKeyValuePairs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "myapp",
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestReader_Read_IgnoresCommentsAndEmptyLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := "# this is a comment\n\nFOO=bar\n# another comment\nBAZ=qux\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestReader_Read_StripsQuotes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := `SECRET="my secret value"` + "\n" + `TOKEN='another token'` + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["SECRET"] != "my secret value" {
		t.Errorf("SECRET: got %q, want %q", got["SECRET"], "my secret value")
	}
	if got["TOKEN"] != "another token" {
		t.Errorf("TOKEN: got %q, want %q", got["TOKEN"], "another token")
	}
}

func TestReader_Read_FileNotExist_ReturnsEmpty(t *testing.T) {
	r := NewReader("/nonexistent/path/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestReader_Read_InvalidFormat_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("INVALID_LINE_NO_EQUALS\n"), 0600); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	r := NewReader(path)
	_, err := r.Read()
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestReader_Read_ValueWithEqualsSign(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Values that contain '=' should be preserved as-is after the first '='
	content := "DATABASE_URL=postgres://user:pass@host/db?sslmode=disable\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("setup: write file: %v", err)
	}

	r := NewReader(path)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "postgres://user:pass@host/db?sslmode=disable"
	if got["DATABASE_URL"] != want {
		t.Errorf("DATABASE_URL: got %q, want %q", got["DATABASE_URL"], want)
	}
}
