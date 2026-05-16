// Package differ provides utilities for comparing environment variable maps.
//
// It supports two-way diffs (Diff) for comparing a source and target environment,
// and three-way diffs (ThreeWay) for merging changes from two branches relative
// to a common base — similar to a standard VCS three-way merge.
//
// Three-way diff results can be inspected for conflicts, resolved automatically
// via PreferLeft/PreferRight policies, or extracted into a flat map with Resolved.
package differ
