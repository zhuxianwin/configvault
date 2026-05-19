package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const maskChar = "***"

func maskValue(value string) string {
	if len(value) == 0 {
		return maskChar
	}
	return maskChar
}

// Fprint writes a human-readable diff summary to w.
// If mask is true, secret values are hidden.
func Fprint(w io.Writer, changes []Change, mask bool) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	// Sort changes by key for deterministic output
	sorted := make([]Change, len(changes))
	copy(sorted, changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		switch c.Type {
		case Added:
			val := c.NewValue
			if mask {
				val = maskValue(val)
			}
			fmt.Fprintf(w, "  + %s=%s\n", c.Key, val)
		case Removed:
			oldVal := c.OldValue
			if mask {
				oldVal = strings.Repeat("*", 3)
			}
			fmt.Fprintf(w, "  - %s=%s\n", c.Key, oldVal)
		case Updated:
			oldVal, newVal := c.OldValue, c.NewValue
			if mask {
				oldVal = maskValue(oldVal)
				newVal = maskValue(newVal)
			}
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Key, oldVal, newVal)
		}
	}
}
