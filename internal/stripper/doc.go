// Package stripper provides functionality to remove environment variable keys
// from a map based on prefix matching, exact key names, or regular expression
// patterns. The original map is never mutated; a new map is returned alongside
// a sorted list of the keys that were removed.
//
// Example:
//
//	res, err := stripper.Strip(env, stripper.Options{
//		Prefixes: []string{"INTERNAL_", "DEBUG_"},
//		Keys:     []string{"LEGACY_FLAG"},
//	})
package stripper
