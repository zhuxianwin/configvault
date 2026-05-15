package diff_test

import (
	"testing"

	"github.com/configvault/internal/diff"
)

func TestCompute_Added(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}

	changes := diff.Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.Added || changes[0].Key != "BAZ" {
		t.Errorf("expected Added change for BAZ, got %+v", changes[0])
	}
}

func TestCompute_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	new := map[string]string{"FOO": "bar"}

	changes := diff.Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.Removed || changes[0].Key != "BAZ" {
		t.Errorf("expected Removed change for BAZ, got %+v", changes[0])
	}
}

func TestCompute_Updated(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "newbar"}

	changes := diff.Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Type != diff.Updated || c.OldValue != "bar" || c.NewValue != "newbar" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestCompute_NoChanges(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes := diff.Compute(m, m)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestHasChanges(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "2"}
	if !diff.HasChanges(old, new) {
		t.Error("expected HasChanges to return true")
	}
	if diff.HasChanges(old, old) {
		t.Error("expected HasChanges to return false for identical maps")
	}
}
