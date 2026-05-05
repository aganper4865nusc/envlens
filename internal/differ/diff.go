package differ

import (
	"fmt"
	"sort"
)

// DiffStatus represents the state of a key when comparing two env files.
type DiffStatus string

const (
	StatusAdded   DiffStatus = "added"
	StatusRemoved DiffStatus = "removed"
	StatusChanged DiffStatus = "changed"
	StatusSame    DiffStatus = "same"
)

// DiffEntry holds the comparison result for a single environment variable key.
type DiffEntry struct {
	Key      string
	Status   DiffStatus
	BaseVal  string
	TargetVal string
}

// Diff compares two parsed env maps (base vs target) and returns a sorted
// slice of DiffEntry describing every key found in either map.
func Diff(base, target map[string]string) []DiffEntry {
	seen := make(map[string]bool)
	var entries []DiffEntry

	for k, bv := range base {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Status: StatusRemoved, BaseVal: bv})
		} else if bv != tv {
			entries = append(entries, DiffEntry{Key: k, Status: StatusChanged, BaseVal: bv, TargetVal: tv})
		} else {
			entries = append(entries, DiffEntry{Key: k, Status: StatusSame, BaseVal: bv, TargetVal: tv})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			entries = append(entries, DiffEntry{Key: k, Status: StatusAdded, TargetVal: tv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// Summary returns a human-readable summary string for a DiffEntry.
func (d DiffEntry) Summary() string {
	switch d.Status {
	case StatusAdded:
		return fmt.Sprintf("+ %s=%q", d.Key, d.TargetVal)
	case StatusRemoved:
		return fmt.Sprintf("- %s=%q", d.Key, d.BaseVal)
	case StatusChanged:
		return fmt.Sprintf("~ %s: %q -> %q", d.Key, d.BaseVal, d.TargetVal)
	default:
		return fmt.Sprintf("  %s=%q", d.Key, d.BaseVal)
	}
}
