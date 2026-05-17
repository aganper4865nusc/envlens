// Package tracker analyses a sequence of environment variable snapshots
// and produces a per-key change history.
//
// Given an ordered slice of env maps (oldest first), Track returns a Result
// that lists every key that changed at least once, along with the old and
// new value for each transition.
//
// Example:
//
//	snaps := []map[string]string{
//		{"DB_HOST": "localhost"},
//		{"DB_HOST": "prod.db.internal"},
//	}
//	result := tracker.Track(snaps)
//	// result.TotalChanges == 1
package tracker
