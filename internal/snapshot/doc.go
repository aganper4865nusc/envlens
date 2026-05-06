// Package snapshot provides functionality to capture, persist, and compare
// point-in-time states of environment variable maps.
//
// A Snapshot records the source label, a UTC timestamp, and a copy of the
// key-value pairs at the time of capture. Snapshots can be serialised to
// JSON files and reloaded later for historical comparison.
//
// Use Take to create a snapshot, Save/Load for persistence, and Compare to
// produce a Delta that describes additions, removals, and value changes
// between two snapshots.
package snapshot
