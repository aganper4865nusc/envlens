// Package profiler provides statistical analysis and comparison of environment
// variable maps. It computes metadata such as total key count, empty value
// count, sensitive key detection, key casing distribution, and unique value
// count. The Compare function allows side-by-side profiling of two environments
// (e.g. dev vs prod) to surface structural differences at a glance.
package profiler
