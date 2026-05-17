// Package tracker records which environment variable keys have changed
// between two snapshots over time, producing a per-key change history.
package tracker

import (
	"sort"
	"time"
)

// Event describes a single change observed for a key.
type Event struct {
	Key       string
	OldValue  string
	NewValue  string
	ChangedAt time.Time
}

// History holds the ordered list of events for a single key.
type History struct {
	Key    string
	Events []Event
}

// Result is the output of a Track call.
type Result struct {
	// Histories contains one entry per key that has at least one event.
	Histories []History
	// TotalChanges is the total number of individual change events recorded.
	TotalChanges int
}

// Track compares a sequence of env snapshots (ordered oldest→newest) and
// builds a per-key change history. Each map represents the full env state
// at that point in time.
func Track(snapshots []map[string]string) Result {
	if len(snapshots) < 2 {
		return Result{}
	}

	// Collect all keys across all snapshots.
	keySet := map[string]struct{}{}
	for _, snap := range snapshots {
		for k := range snap {
			keySet[k] = struct{}{}
		}
	}

	now := time.Now().UTC()
	historyMap := map[string]*History{}
	total := 0

	for key := range keySet {
		var events []Event
		for i := 1; i < len(snapshots); i++ {
			prev := snapshots[i-1][key]
			curr := snapshots[i][key]
			if prev != curr {
				events = append(events, Event{
					Key:       key,
					OldValue:  prev,
					NewValue:  curr,
					ChangedAt: now,
				})
				total++
			}
		}
		if len(events) > 0 {
			historyMap[key] = &History{Key: key, Events: events}
		}
	}

	keys := make([]string, 0, len(historyMap))
	for k := range historyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	histories := make([]History, 0, len(keys))
	for _, k := range keys {
		histories = append(histories, *historyMap[k])
	}

	return Result{Histories: histories, TotalChanges: total}
}

// KeysChanged returns the sorted list of keys that have at least one change event.
func KeysChanged(r Result) []string {
	out := make([]string, 0, len(r.Histories))
	for _, h := range r.Histories {
		out = append(out, h.Key)
	}
	return out
}
