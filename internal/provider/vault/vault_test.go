package vault_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"configvault/internal/provider/vault"
)

func TestProvider_GetSecrets_ParsesResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"data":{"DB_HOST":"localhost","API_KEY":"abc123"}}}`))  
	}))
	defer ts.Close()

	p := vault.New(ts.URL, "test-token", nil)
	secrets, err := p.GetSecrets(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}

func TestProvider_GetSecrets_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":["path not found"]}`)) 
	}))
	defer ts.Close()

	p := vault.New(ts.URL, "any-token", nil)
	_, err := p.GetSecrets(context.Background(), "secret/data/missing")
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestProvider_Name(t *testing.T) {
	p := vault.New("http://localhost:8200", "token", nil)
	if p.Name() != "vault" {
		t.Errorf("expected name vault, got %q", p.Name())
	}
}
