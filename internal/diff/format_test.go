package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestFprint_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, []Change{}, false)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %q", buf.String())
	}
}

func TestFprint_Added_MasksValue(t *testing.T) {
	changes := []Change{
		{Key: "SECRET", Type: Added, NewValue: "supersecret"},
	}
	var buf bytes.Buffer
	Fprint(&buf, changes, true)
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected value to be masked, got: %q", out)
	}
	if !strings.Contains(out, "+ SECRET") {
		t.Errorf("expected added marker, got: %q", out)
	}
}

func TestFprint_Added_ShowsValue(t *testing.T) {
	changes := []Change{
		{Key: "API_KEY", Type: Added, NewValue: "abc123"},
	}
	var buf bytes.Buffer
	Fprint(&buf, changes, false)
	out := buf.String()
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected plain value, got: %q", out)
	}
}

func TestFprint_Removed(t *testing.T) {
	changes := []Change{
		{Key: "OLD_KEY", Type: Removed, OldValue: "oldval"},
	}
	var buf bytes.Buffer
	Fprint(&buf, changes, false)
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removed marker, got: %q", out)
	}
	if !strings.Contains(out, "oldval") {
		t.Errorf("expected old value in output, got: %q", out)
	}
}

func TestFprint_Updated_SortedOutput(t *testing.T) {
	changes := []Change{
		{Key: "Z_KEY", Type: Updated, OldValue: "old2", NewValue: "new2"},
		{Key: "A_KEY", Type: Updated, OldValue: "old1", NewValue: "new1"},
	}
	var buf bytes.Buffer
	Fprint(&buf, changes, false)
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if aIdx == -1 || zIdx == -1 {
		t.Fatalf("expected both keys in output, got: %q", out)
	}
	if aIdx > zIdx {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}
