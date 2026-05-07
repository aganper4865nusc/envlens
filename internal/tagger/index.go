package tagger

// Index builds an inverted index mapping each tag to the list of keys
// that carry that tag. The key lists are sorted alphabetically.
// This is useful for quickly querying "which keys are tagged X?".
func Index(results []Result) map[string][]string {
	idx := map[string][]string{}
	for _, r := range results {
		for _, tag := range r.Tags {
			idx[tag] = append(idx[tag], r.Key)
		}
	}
	for tag := range idx {
		sortStrings(idx[tag])
	}
	return idx
}

// TagsFor returns the tags assigned to a specific key from a result slice.
// Returns nil if the key is not found.
func TagsFor(key string, results []Result) []string {
	for _, r := range results {
		if r.Key == key {
			return r.Tags
		}
	}
	return nil
}
