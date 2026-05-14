package scoper

// MergeOptions controls how scoped results are merged back into a flat map.
type MergeOptions struct {
	// Prefix, if set, is prepended to every key in the merged output.
	Prefix string
	// Overwrite allows later results to overwrite earlier ones on key conflict.
	Overwrite bool
}

// MergeResult is returned by Merge.
type MergeResult struct {
	Env       map[string]string
	Conflicts []string
}

// Merge combines multiple scoped Results into a single flat map.
// Conflicts are recorded when Overwrite is false and a key collision occurs.
func Merge(results []Result, opts MergeOptions) MergeResult {
	out := make(map[string]string)
	var conflicts []string

	for _, r := range results {
		for k, v := range r.Env {
			key := opts.Prefix + k
			if _, exists := out[key]; exists {
				if !opts.Overwrite {
					conflicts = append(conflicts, key)
					continue
				}
			}
			out[key] = v
		}
	}

	return MergeResult{
		Env:       out,
		Conflicts: conflicts,
	}
}
