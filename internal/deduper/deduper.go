// Package deduper removes duplicate environment variable entries,
// keeping the last-defined value for each key (consistent with shell semantics).
package deduper

// Result holds the output of a deduplication pass.
type Result struct {
	// Env is the deduplicated map of key → value.
	Env map[string]string
	// Duplicates maps each key that appeared more than once to the number of
	// times it was seen (including the winning entry).
	Duplicates map[string]int
}

// Dedupe processes an ordered slice of key/value pairs and returns a Result.
// Pairs are supplied as []Entry so that order — and therefore "last wins"
// semantics — is preserved even though a plain map has no ordering guarantee.
type Entry struct {
	Key   string
	Value string
}

// Dedupe deduplicates the provided entries.
// The last entry for a given key wins.
func Dedupe(entries []Entry) Result {
	seen := make(map[string]int, len(entries))
	env := make(map[string]string, len(entries))

	for _, e := range entries {
		seen[e.Key]++
		env[e.Key] = e.Value
	}

	duplicates := make(map[string]int)
	for k, count := range seen {
		if count > 1 {
			duplicates[k] = count
		}
	}

	return Result{
		Env:        env,
		Duplicates: duplicates,
	}
}

// FromMap is a convenience wrapper when the caller already has a plain map and
// no ordering concerns — it returns a Result with no duplicates.
func FromMap(m map[string]string) Result {
	env := make(map[string]string, len(m))
	for k, v := range m {
		env[k] = v
	}
	return Result{
		Env:        env,
		Duplicates: make(map[string]int),
	}
}
