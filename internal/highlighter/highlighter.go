// Package highlighter annotates environment variable maps with change markers
// relative to a baseline, making it easy to surface what is new, removed, or
// modified when comparing two env snapshots or files.
package highlighter

import "sort"

// Status describes the change state of a single key.
type Status string

const (
	Added    Status = "added"
	Removed  Status = "removed"
	Changed  Status = "changed"
	Unchanged Status = "unchanged"
)

// Entry holds a key, its resolved value, and its change status.
type Entry struct {
	Key      string
	Value    string
	OldValue string // populated only when Status == Changed
	Status   Status
}

// Result is the full output of a Highlight call.
type Result struct {
	Entries  []Entry
	Added    int
	Removed  int
	Changed  int
	Unchanged int
}

// Highlight compares current against baseline and returns annotated entries
// sorted alphabetically by key.
func Highlight(baseline, current map[string]string) Result {
	keys := unionKeys(baseline, current)
	sort.Strings(keys)

	var res Result
	for _, k := range keys {
		bVal, inBase := baseline[k]
		cVal, inCurr := current[k]

		switch {
		case inCurr && !inBase:
			res.Entries = append(res.Entries, Entry{Key: k, Value: cVal, Status: Added})
			res.Added++
		case inBase && !inCurr:
			res.Entries = append(res.Entries, Entry{Key: k, Value: bVal, Status: Removed})
			res.Removed++
		case bVal != cVal:
			res.Entries = append(res.Entries, Entry{Key: k, Value: cVal, OldValue: bVal, Status: Changed})
			res.Changed++
		default:
			res.Entries = append(res.Entries, Entry{Key: k, Value: cVal, Status: Unchanged})
			res.Unchanged++
		}
	}
	return res
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
