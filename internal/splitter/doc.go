// Package splitter partitions an environment variable map into named
// buckets using prefix or regex pattern rules.
//
// Rules are evaluated in order; the first matching rule wins. Keys that
// do not match any rule are placed in the Unmatched map of the Result.
//
// Example:
//
//	res, err := splitter.Split(env, []splitter.Rule{
//		{Name: "database", Prefix: "DB_"},
//		{Name: "auth",     Pattern: `^(JWT|OAUTH)_`},
//	})
package splitter
