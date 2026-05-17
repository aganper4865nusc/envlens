// Package selector provides fine-grained selection of environment variable keys
// from a map using one or more criteria: explicit key names, key prefixes, or
// regular expression patterns.
//
// Selection can be inverted to exclude matching keys instead of including them.
// Results include sorted slices of selected and skipped keys for deterministic
// downstream processing.
//
// Example:
//
//	res, err := selector.Select(env, selector.Options{
//		Prefixes: []string{"DB_", "REDIS_"},
//		Invert:   false,
//	})
package selector
