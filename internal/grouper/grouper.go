// Package grouper organizes environment variables into named groups
// based on key prefix conventions (e.g., DB_, AWS_, APP_).
package grouper

import (
	"sort"
	"strings"
)

// Group holds a named collection of env vars sharing a common prefix.
type Group struct {
	Name string
	Keys []string
	Env  map[string]string
}

// Result is the output of a Group operation.
type Result struct {
	Groups   []Group
	Ungrouped map[string]string
}

// GroupBy partitions the env map into groups by the provided prefixes.
// Keys that match no prefix are placed in Ungrouped.
// Prefixes are matched case-insensitively against the key.
func GroupBy(env map[string]string, prefixes []string) Result {
	groupMap := make(map[string]*Group, len(prefixes))
	for _, p := range prefixes {
		upper := strings.ToUpper(p)
		groupMap[upper] = &Group{
			Name: upper,
			Env:  make(map[string]string),
		}
	}

	ungrouped := make(map[string]string)

	for k, v := range env {
		matched := false
		for _, p := range prefixes {
			upper := strings.ToUpper(p)
			if strings.HasPrefix(strings.ToUpper(k), upper) {
				groupMap[upper].Env[k] = v
				groupMap[upper].Keys = append(groupMap[upper].Keys, k)
				matched = true
				break
			}
		}
		if !matched {
			ungrouped[k] = v
		}
	}

	groups := make([]Group, 0, len(prefixes))
	for _, p := range prefixes {
		upper := strings.ToUpper(p)
		g := groupMap[upper]
		sort.Strings(g.Keys)
		groups = append(groups, *g)
	}

	return Result{
		Groups:    groups,
		Ungrouped: ungrouped,
	}
}
