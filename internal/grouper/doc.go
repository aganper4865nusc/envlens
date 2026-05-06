// Package grouper partitions an environment variable map into named groups
// based on key prefix conventions.
//
// This is useful for visualising which subsystem (database, cloud provider,
// application) each variable belongs to, and for targeted diffing or auditing
// of a specific prefix namespace.
package grouper
