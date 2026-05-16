package differ

// ThreeWayResult holds the result of a three-way diff between a base,
// left (ours), and right (theirs) environment map.
type ThreeWayResult struct {
	Key       string
	Base      string
	Left      string
	Right     string
	Conflict  bool
	Resolution string
}

// ThreeWayOptions controls how conflicts are resolved.
type ThreeWayOptions struct {
	// PreferLeft resolves conflicts by picking the left (ours) value.
	// If false, conflicts are left unresolved (Conflict=true, Resolution="").
	PreferLeft bool
	// PreferRight resolves conflicts by picking the right (theirs) value.
	// Takes precedence over PreferLeft.
	PreferRight bool
}

// ThreeWay performs a three-way diff of environment maps.
// base is the common ancestor; left is "ours"; right is "theirs".
func ThreeWay(base, left, right map[string]string, opts ThreeWayOptions) []ThreeWayResult {
	keys := unionAll(base, left, right)
	results := make([]ThreeWayResult, 0, len(keys))

	for _, k := range keys {
		bv := base[k]
		lv := left[k]
		rv := right[k]

		r := ThreeWayResult{Key: k, Base: bv, Left: lv, Right: rv}

		switch {
		case lv == rv:
			// Both sides agree — no conflict.
			r.Resolution = lv
		case lv == bv:
			// Only right changed.
			r.Resolution = rv
		case rv == bv:
			// Only left changed.
			r.Resolution = lv
		default:
			// Both sides changed differently — conflict.
			r.Conflict = true
			if opts.PreferRight {
				r.Resolution = rv
			} else if opts.PreferLeft {
				r.Resolution = lv
			}
		}

		results = append(results, r)
	}

	return results
}

// Resolved returns a flat map of resolved values, omitting unresolved conflicts.
func Resolved(results []ThreeWayResult) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		if !r.Conflict || r.Resolution != "" {
			out[r.Key] = r.Resolution
		}
	}
	return out
}

// Conflicts returns only the conflicting entries.
func Conflicts(results []ThreeWayResult) []ThreeWayResult {
	var out []ThreeWayResult
	for _, r := range results {
		if r.Conflict {
			out = append(out, r)
		}
	}
	return out
}

func unionAll(maps ...map[string]string) []string {
	seen := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
