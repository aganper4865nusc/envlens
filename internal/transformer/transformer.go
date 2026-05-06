// Package transformer provides utilities for transforming environment variable
// maps by applying key/value mutations such as prefix injection, key renaming,
// and value uppercasing.
package transformer

import "strings"

// Options configures the transformation behaviour.
type Options struct {
	// AddPrefix prepends a string to every key.
	AddPrefix string
	// StripPrefix removes a leading string from every key.
	StripPrefix string
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// UppercaseValues converts all values to uppercase.
	UppercaseValues bool
}

// Result holds the transformed map and a log of changes made.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Change records a single key or value mutation.
type Change struct {
	Key      string
	OldKey   string // non-empty when the key itself was renamed
	OldValue string
	NewValue string
}

// Transform applies the given Options to src and returns a Result.
// src is never modified.
func Transform(src map[string]string, opts Options) Result {
	out := make(map[string]string, len(src))
	var changes []Change

	for k, v := range src {
		newKey := k
		newVal := v

		if opts.StripPrefix != "" && strings.HasPrefix(newKey, opts.StripPrefix) {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}
		if opts.AddPrefix != "" {
			newKey = opts.AddPrefix + newKey
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.UppercaseValues {
			newVal = strings.ToUpper(v)
		}

		if newKey != k || newVal != v {
			c := Change{Key: newKey, OldValue: v, NewValue: newVal}
			if newKey != k {
				c.OldKey = k
			}
			changes = append(changes, c)
		}
		out[newKey] = newVal
	}

	return Result{Env: out, Changes: changes}
}
