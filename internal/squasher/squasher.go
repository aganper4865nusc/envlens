// Package squasher merges multiple env maps into a single flat map,
// applying a defined conflict resolution strategy and tracking which
// keys were overwritten and by which source.
package squasher

import "sort"

// Strategy controls how conflicts are resolved when the same key
// appears in more than one source.
type Strategy int

const (
	// FirstWins keeps the value from the earliest source.
	FirstWins Strategy = iota
	// LastWins keeps the value from the latest source (default merge behaviour).
	LastWins
)

// Conflict records a key that existed in more than one source.
type Conflict struct {
	Key      string
	Kept     string // value that was retained
	Dropped  string // value that was discarded
	WonIndex int    // zero-based index of the winning source
}

// Result is returned by Squash.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Squash merges sources in order according to the given strategy.
// Sources later in the slice have higher priority when Strategy is
// LastWins (the default); earlier sources win when Strategy is FirstWins.
func Squash(sources []map[string]string, strategy Strategy) Result {
	out := make(map[string]string)
	origin := make(map[string]int) // key -> source index that currently owns the value
	var conflicts []Conflict

	for idx, src := range sources {
		for k, v := range src {
			if existing, exists := out[k]; exists {
				switch strategy {
				case FirstWins:
					conflicts = append(conflicts, Conflict{
						Key:      k,
						Kept:     existing,
						Dropped:  v,
						WonIndex: origin[k],
					})
				case LastWins:
					conflicts = append(conflicts, Conflict{
						Key:      k,
						Kept:     v,
						Dropped:  existing,
						WonIndex: idx,
					})
					out[k] = v
					origin[k] = idx
				}
			} else {
				out[k] = v
				origin[k] = idx
			}
		}
	}

	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Env: out, Conflicts: conflicts}
}
