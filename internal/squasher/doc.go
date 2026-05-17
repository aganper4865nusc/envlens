// Package squasher combines multiple environment variable maps into a
// single flat map. Unlike merger.Merge, squasher makes conflict
// resolution explicit via a Strategy and surfaces every overwrite as
// a Conflict in the returned Result, making it straightforward for
// callers to audit exactly which values were dropped and from which
// source layer.
//
// Supported strategies:
//
//	FirstWins – the first source to define a key owns its value.
//	LastWins  – the last source to define a key wins (mirrors typical
//	            shell export semantics).
package squasher
