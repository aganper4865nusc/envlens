// Package normalizer provides utilities to normalize environment variable
// maps by standardizing key casing, trimming whitespace, and collapsing
// redundant entries into a canonical form.
package normalizer

import (
	"strings"
)

// Options controls how normalization is applied.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// TrimKeys strips leading and trailing whitespace from keys.
	TrimKeys bool
	// RemoveEmpty drops entries whose value is empty after trimming.
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys: true,
		TrimValues:    true,
		TrimKeys:      true,
		RemoveEmpty:   false,
	}
}

// Result holds the normalized map and metadata about changes made.
type Result struct {
	// Env is the normalized environment map.
	Env map[string]string
	// RenamedKeys lists keys that were renamed (e.g. due to casing).
	RenamedKeys []string
	// TrimmedValues lists keys whose values were trimmed.
	TrimmedValues []string
	// DroppedKeys lists keys that were removed because their value was empty.
	DroppedKeys []string
}

// Normalize applies the given Options to the input env map and returns a
// Result containing the normalized copy and a summary of changes.
func Normalize(env map[string]string, opts Options) Result {
	result := Result{
		Env: make(map[string]string, len(env)),
	}

	for k, v := range env {
		newKey := k
		if opts.TrimKeys {
			newKey = strings.TrimSpace(newKey)
		}
		if opts.UppercaseKeys {
			upper := strings.ToUpper(newKey)
			if upper != newKey {
				result.RenamedKeys = append(result.RenamedKeys, k)
			}
			newKey = upper
		}

		newVal := v
		if opts.TrimValues {
			trimmed := strings.TrimSpace(v)
			if trimmed != v {
				result.TrimmedValues = append(result.TrimmedValues, newKey)
			}
			newVal = trimmed
		}

		if opts.RemoveEmpty && newVal == "" {
			result.DroppedKeys = append(result.DroppedKeys, newKey)
			continue
		}

		result.Env[newKey] = newVal
	}

	return result
}
