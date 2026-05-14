// Package scoper partitions an environment variable map into named scopes
// based on key-prefix conventions (e.g. PROD_, STAGING_).
//
// # Basic usage
//
//	rule := scoper.Rule{Scope: "prod", Prefix: "PROD_", StripPrefix: true}
//	result := scoper.Scope(env, rule)
//	// result.Env contains only PROD_* keys, with the prefix removed.
//
// # Multiple scopes
//
//	rules := []scoper.Rule{
//	    {Scope: "prod",    Prefix: "PROD_",    StripPrefix: true},
//	    {Scope: "staging", Prefix: "STAGING_", StripPrefix: true},
//	}
//	results := scoper.ScopeAll(env, rules)
//
// # Merging scopes
//
//	merged := scoper.Merge(results, scoper.MergeOptions{Overwrite: false})
package scoper
