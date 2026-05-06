// Package patcher provides functionality to apply declarative patch operations
// to an environment variable map.
//
// Supported operations:
//
//	set KEY=VALUE   — add or overwrite a key
//	unset KEY       — remove a key
//	rename OLD NEW  — rename a key, preserving its value
//
// Patch files can be loaded with LoadOps and applied with Patch.
// The original map is never mutated; a new map is always returned.
package patcher
