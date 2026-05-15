// Package zipper merges two env maps key-by-key using a user-supplied
// combiner function, producing a new map of transformed values.
package zipper

import "sort"

// PairResult holds the outcome of combining a single key from two env maps.
type PairResult struct {
	Key        string
	LeftValue  string
	RightValue string
	Result     string
	LeftOnly   bool
	RightOnly  bool
}

// Options controls Zip behaviour.
type Options struct {
	// IncludeLeftOnly includes keys present only in the left map.
	IncludeLeftOnly bool
	// IncludeRightOnly includes keys present only in the right map.
	IncludeRightOnly bool
}

// DefaultOptions returns sensible defaults (include all keys).
func DefaultOptions() Options {
	return Options{IncludeLeftOnly: true, IncludeRightOnly: true}
}

// CombineFunc receives the left and right values for a key and returns the
// combined result. Either value may be empty when the key is absent on that
// side.
type CombineFunc func(left, right string) string

// Zip iterates over the union of keys in left and right, calls combine for
// each pair, and returns the resulting slice ordered by key.
func Zip(left, right map[string]string, opts Options, combine CombineFunc) []PairResult {
	keys := unionKeys(left, right)
	sort.Strings(keys)

	var results []PairResult
	for _, k := range keys {
		lv, inLeft := left[k]
		rv, inRight := right[k]

		if inLeft && !inRight && !opts.IncludeLeftOnly {
			continue
		}
		if inRight && !inLeft && !opts.IncludeRightOnly {
			continue
		}

		results = append(results, PairResult{
			Key:        k,
			LeftValue:  lv,
			RightValue: rv,
			Result:     combine(lv, rv),
			LeftOnly:   inLeft && !inRight,
			RightOnly:  inRight && !inLeft,
		})
	}
	return results
}

// ToMap converts a Zip result slice into a plain map[string]string.
func ToMap(pairs []PairResult) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		out[p.Key] = p.Result
	}
	return out
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
