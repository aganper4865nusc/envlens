// Package walker provides directory traversal for discovering .env files
// within a project tree.
//
// It supports configurable filename patterns, maximum recursion depth, and
// optional skipping of hidden directories. Results are returned sorted by
// path for deterministic output.
//
// Example:
//
//	results, err := walker.Walk(".", walker.Options{
//		MaxDepth:   3,
//		SkipHidden: true,
//	})
package walker
