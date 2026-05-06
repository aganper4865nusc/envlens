// Package renamer provides utilities for bulk-renaming environment variable keys
// according to a set of rename rules, tracking all applied and skipped renames.
package renamer

import "strings"

// Rule describes a single rename operation.
type Rule struct {
	From string // original key name
	To   string // desired key name
}

// Result holds the outcome of a Rename operation.
type Result struct {
	Env     map[string]string // resulting environment map
	Applied []Rule            // rules that were successfully applied
	Skipped []Rule            // rules whose From key was not found
	Conflicts []Rule          // rules whose To key already existed (not applied)
}

// Rename applies the given rules to env, returning a new map and a Result
// describing what happened. The original map is never modified.
//
// Conflict policy: if the destination key already exists in env (and is not
// itself being renamed away in the same batch), the rule is recorded as a
// conflict and skipped.
func Rename(env map[string]string, rules []Rule) Result {
	// Build a set of keys that will be vacated by this batch so we can
	// distinguish a real conflict from a swap within the same rule set.
	vacated := make(map[string]bool, len(rules))
	for _, r := range rules {
		if _, ok := env[r.From]; ok {
			vacated[r.From] = true
		}
	}

	out := copyMap(env)
	result := Result{Env: out}

	for _, r := range rules {
		r.From = strings.TrimSpace(r.From)
		r.To = strings.TrimSpace(r.To)

		val, exists := env[r.From]
		if !exists {
			result.Skipped = append(result.Skipped, r)
			continue
		}

		// Conflict: destination already occupied by a key that is NOT being vacated.
		if _, occupied := env[r.To]; occupied && !vacated[r.To] {
			result.Conflicts = append(result.Conflicts, r)
			continue
		}

		delete(out, r.From)
		out[r.To] = val
		result.Applied = append(result.Applied, r)
	}

	return result
}

func copyMap(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
