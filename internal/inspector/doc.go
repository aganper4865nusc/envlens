// Package inspector provides a unified analysis pass over a single env map.
//
// It combines the output of the profiler, linter, auditor, and scorer packages
// into a single Report value, making it easy to display a complete picture of
// an environment file's health without wiring each subsystem individually.
//
// Usage:
//
//	env, _ := parser.ParseFile("production.env")
//	report := inspector.Inspect("production.env", env)
//	inspector.WriteText(os.Stdout, report)
package inspector
