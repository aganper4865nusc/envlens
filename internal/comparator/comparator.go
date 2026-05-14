// Package comparator provides side-by-side comparison of two env maps,
// producing a structured result with match rate and per-key status.
package comparator

import "sort"

// Status represents the comparison result for a single key.
type Status string

const (
	StatusMatch   Status = "match"
	StatusMissing Status = "missing" // present in left, absent in right
	StatusExtra   Status = "extra"   // absent in left, present in right
	StatusDiffer  Status = "differ"
)

// Entry holds the comparison result for one key.
type Entry struct {
	Key        string
	LeftValue  string
	RightValue string
	Status     Status
}

// Result is the full output of a comparison.
type Result struct {
	Entries   []Entry
	MatchRate float64 // 0.0–1.0
	Total     int
	Matched   int
}

// Compare compares two env maps (left vs right) and returns a Result.
// Keys are evaluated from the union of both maps.
func Compare(left, right map[string]string) Result {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	var entries []Entry
	matched := 0

	for _, k := range keys {
		lv, linOK := left[k]
		rv, rinOK := right[k]

		var status Status
		switch {
		case linOK && !rinOK:
			status = StatusMissing
		case !linOK && rinOK:
			status = StatusExtra
		case lv == rv:
			status = StatusMatch
			matched++
		default:
			status = StatusDiffer
		}

		entries = append(entries, Entry{
			Key:        k,
			LeftValue:  lv,
			RightValue: rv,
			Status:     status,
		})
	}

	total := len(entries)
	var rate float64
	if total > 0 {
		rate = float64(matched) / float64(total)
	}

	return Result{
		Entries:   entries,
		MatchRate: rate,
		Total:     total,
		Matched:   matched,
	}
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
