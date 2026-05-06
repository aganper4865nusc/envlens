package merger

// MergeResult holds the result of merging multiple env maps.
type MergeResult struct {
	// Merged is the final merged map (later sources override earlier ones).
	Merged map[string]string
	// Overrides maps each key to the list of source indices that defined it.
	Overrides map[string][]int
}

// Merge combines multiple env variable maps into one.
// Sources are applied in order; later sources override earlier ones.
// Keys that appear in more than one source are tracked in Overrides.
func Merge(sources ...map[string]string) MergeResult {
	merged := make(map[string]string)
	overrides := make(map[string][]int)

	for idx, src := range sources {
		for k, v := range src {
			if _, exists := merged[k]; exists {
				// Record that this key was already present before this source.
				overrides[k] = append(overrides[k], idx)
			} else {
				// First time we see this key — record the originating source.
				overrides[k] = []int{idx}
			}
			merged[k] = v
		}
	}

	// Clean up keys that were only seen in a single source (not real overrides).
	for k, indices := range overrides {
		if len(indices) < 2 {
			delete(overrides, k)
		}
	}

	return MergeResult{
		Merged:    merged,
		Overrides: overrides,
	}
}
