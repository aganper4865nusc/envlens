// Package sorter provides utilities for sorting and grouping environment
// variable maps by key prefix, alphabetical order, or custom criteria.
package sorter

import (
	"sort"
	"strings"
)

// SortMode controls how keys are ordered.
type SortMode int

const (
	// Alphabetical sorts keys in ascending lexicographic order.
	Alphabetical SortMode = iota
	// AlphabeticalDesc sorts keys in descending lexicographic order.
	AlphabeticalDesc
	// ByPrefix groups keys by their prefix (segment before the first '_').
	ByPrefix
)

// Result holds sorted keys and, when ByPrefix mode is used, a grouping map.
type Result struct {
	// Keys contains all keys in sorted order.
	Keys []string
	// Groups maps each prefix to its member keys (populated in ByPrefix mode).
	Groups map[string][]string
}

// Sort orders the keys of env according to the given mode and returns a Result.
func Sort(env map[string]string, mode SortMode) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch mode {
	case AlphabeticalDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case ByPrefix:
		sort.Strings(keys)
		groups := groupByPrefix(keys)
		return Result{Keys: keys, Groups: groups}
	default:
		sort.Strings(keys)
	}

	return Result{Keys: keys, Groups: nil}
}

// groupByPrefix partitions a sorted slice of keys by the segment that precedes
// the first underscore. Keys without an underscore are placed under "".
func groupByPrefix(keys []string) map[string][]string {
	groups := make(map[string][]string)
	for _, k := range keys {
		prefix := ""
		if idx := strings.Index(k, "_"); idx > 0 {
			prefix = k[:idx]
		}
		groups[prefix] = append(groups[prefix], k)
	}
	return groups
}
