// Package labeler provides rule-based label assignment for environment variable keys.
//
// Rules may match keys by prefix, regular expression, or explicit name list.
// Multiple rules may apply to the same key, resulting in multiple labels.
// Labels are sorted and deduplicated per key.
//
// Example usage:
//
//	rules := []labeler.Rule{
//		{Label: "database", Prefix: "DB_"},
//		{Label: "secret",   Pattern: "(?i)secret|password|token"},
//	}
//	result, err := labeler.Label(env, rules)
package labeler
