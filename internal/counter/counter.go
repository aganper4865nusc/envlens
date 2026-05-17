// Package counter provides utilities for counting and summarizing
// key-value statistics across one or more env maps.
package counter

import "sort"

// KeyStat holds aggregate statistics for a single key across multiple env maps.
type KeyStat struct {
	Key        string
	Occurrence int
	EmptyCount int
	UniqueVals int
}

// Result is the output of a Count operation.
type Result struct {
	Stats      []KeyStat
	TotalKeys  int
	TotalEmpty int
}

// Count aggregates per-key statistics across multiple env maps.
// Keys are sorted alphabetically in the result.
func Count(sources []map[string]string) Result {
	type agg struct {
		occurrence int
		emptyCount int
		vals       map[string]struct{}
	}

	aggs := make(map[string]*agg)

	for _, src := range sources {
		for k, v := range src {
			a, ok := aggs[k]
			if !ok {
				a = &agg{vals: make(map[string]struct{})}
				aggs[k] = a
			}
			a.occurrence++
			if v == "" {
				a.emptyCount++
			} else {
				a.vals[v] = struct{}{}
			}
		}
	}

	keys := make([]string, 0, len(aggs))
	for k := range aggs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	stats := make([]KeyStat, 0, len(keys))
	totalEmpty := 0
	for _, k := range keys {
		a := aggs[k]
		totalEmpty += a.emptyCount
		stats = append(stats, KeyStat{
			Key:        k,
			Occurrence: a.occurrence,
			EmptyCount: a.emptyCount,
			UniqueVals: len(a.vals),
		})
	}

	return Result{
		Stats:      stats,
		TotalKeys:  len(keys),
		TotalEmpty: totalEmpty,
	}
}
