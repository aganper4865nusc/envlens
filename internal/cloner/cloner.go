// Package cloner provides utilities for deep-copying environment variable maps
// with optional key/value transformations applied during the clone operation.
package cloner

import (
	"strings"
)

// Options controls how the clone is performed.
type Options struct {
	// UppercaseKeys normalises all keys to uppercase during cloning.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// ExcludeKeys is a set of keys to omit from the cloned map.
	ExcludeKeys []string
}

// DefaultOptions returns a zero-value Options (no transformations).
func DefaultOptions() Options {
	return Options{}
}

// Result is returned by Clone and carries the cloned environment along
// with metadata about what was excluded or transformed.
type Result struct {
	// Env is the cloned environment map.
	Env map[string]string
	// ExcludedKeys lists any keys that were dropped due to ExcludeKeys.
	ExcludedKeys []string
	// TransformedKeys maps original key → new key for every key that was
	// renamed by UppercaseKeys.
	TransformedKeys map[string]string
}

// Clone deep-copies src applying the transformations described by opts.
func Clone(src map[string]string, opts Options) Result {
	excludeSet := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = struct{}{}
	}

	result := Result{
		Env:             make(map[string]string, len(src)),
		TransformedKeys: make(map[string]string),
	}

	for k, v := range src {
		if _, skip := excludeSet[k]; skip {
			result.ExcludedKeys = append(result.ExcludedKeys, k)
			continue
		}

		newKey := k
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(k)
			if newKey != k {
				result.TransformedKeys[k] = newKey
			}
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
		}

		result.Env[newKey] = newVal
	}

	return result
}
