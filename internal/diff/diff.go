package diff

// Change represents a single change between two secret maps.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// ChangeType describes the kind of change.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Updated ChangeType = "updated"
)

// Compute returns the list of changes between the old and new secret maps.
func Compute(old, new map[string]string) []Change {
	var changes []Change

	for k, newVal := range new {
		oldVal, exists := old[k]
		if !exists {
			changes = append(changes, Change{Key: k, NewValue: newVal, Type: Added})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Type: Updated})
		}
	}

	for k, oldVal := range old {
		if _, exists := new[k]; !exists {
			changes = append(changes, Change{Key: k, OldValue: oldVal, Type: Removed})
		}
	}

	return changes
}

// HasChanges returns true if there are any differences between old and new.
func HasChanges(old, new map[string]string) bool {
	return len(Compute(old, new)) > 0
}
