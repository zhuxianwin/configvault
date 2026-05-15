package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/configvault/internal/diff"
)

func TestFprint_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	diff.Fprint(&buf, nil, false)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", buf.String())
	}
}

func TestFprint_Added_MasksValue(t *testing.T) {
	changes := []diff.Change{
		{Key: "SECRET", NewValue: "supersecret", Type: diff.Added},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, changes, false)
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Error("expected value to be masked")
	}
	if !strings.Contains(out, "+ SECRET") {
		t.Errorf("expected '+ SECRET' in output, got: %s", out)
	}
}

func TestFprint_Added_ShowsValue(t *testing.T) {
	changes := []diff.Change{
		{Key: "SECRET", NewValue: "supersecret", Type: diff.Added},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, changes, true)
	out := buf.String()
	if !strings.Contains(out, "supersecret") {
		t.Errorf("expected plain value in output, got: %s", out)
	}
}

func TestFprint_Removed(t *testing.T) {
	changes := []diff.Change{
		{Key: "OLD_KEY", OldValue: "val", Type: diff.Removed},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, changes, false)
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected '- OLD_KEY' in output, got: %s", out)
	}
}

func TestFprint_Updated_SortedOutput(t *testing.T) {
	changes := []diff.Change{
		{Key: "Z_KEY", OldValue: "a", NewValue: "b", Type: diff.Updated},
		{Key: "A_KEY", OldValue: "x", NewValue: "y", Type: diff.Updated},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, changes, false)
	out := buf.String()
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected A_KEY to appear before Z_KEY in sorted output")
	}
}
