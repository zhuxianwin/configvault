package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// maskValue replaces a secret value with asterisks.
func maskValue(v string) string {
	if len(v) == 0 {
		return ""
	}
	return strings.Repeat("*", 8)
}

// Fprint writes a human-readable summary of changes to w.
// Secret values are masked unless showValues is true.
func Fprint(w io.Writer, changes []Change, showValues bool) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		switch c.Type {
		case Added:
			val := maskValue(c.NewValue)
			if showValues {
				val = c.NewValue
			}
			fmt.Fprintf(w, "  + %s=%s\n", c.Key, val)
		case Removed:
			fmt.Fprintf(w, "  - %s\n", c.Key)
		case Updated:
			oldVal, newVal := maskValue(c.OldValue), maskValue(c.NewValue)
			if showValues {
				oldVal, newVal = c.OldValue, c.NewValue
			}
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Key, oldVal, newVal)
		}
	}
}
