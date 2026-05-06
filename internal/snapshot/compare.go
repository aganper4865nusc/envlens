package snapshot

import "sort"

// Delta describes the difference between two snapshots.
type Delta struct {
	Added   map[string]string // keys present in newer but not older
	Removed map[string]string // keys present in older but not newer
	Changed map[string][2]string // key -> [oldVal, newVal]
	Unchanged []string          // keys with identical values
}

// Compare returns a Delta between an older and a newer Snapshot.
func Compare(older, newer *Snapshot) *Delta {
	d := &Delta{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, newVal := range newer.Env {
		oldVal, exists := older.Env[k]
		if !exists {
			d.Added[k] = newVal
		} else if oldVal != newVal {
			d.Changed[k] = [2]string{oldVal, newVal}
		} else {
			d.Unchanged = append(d.Unchanged, k)
		}
	}

	for k, oldVal := range older.Env {
		if _, exists := newer.Env[k]; !exists {
			d.Removed[k] = oldVal
		}
	}

	sort.Strings(d.Unchanged)
	return d
}

// HasChanges returns true if there are any additions, removals, or changes.
func (d *Delta) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
