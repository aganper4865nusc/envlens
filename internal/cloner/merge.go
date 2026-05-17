package cloner

// MergePolicy controls what happens when a key exists in both the base
// environment and the cloned overlay.
type MergePolicy int

const (
	// PolicyOverwrite lets the overlay value replace the base value (default).
	PolicyOverwrite MergePolicy = iota
	// PolicyKeepBase retains the base value when a conflict is detected.
	PolicyKeepBase
)

// MergeResult is the output of CloneInto.
type MergeResult struct {
	// Env is the merged environment.
	Env map[string]string
	// Conflicts lists keys where both base and overlay had different values.
	Conflicts []string
}

// CloneInto clones src (applying opts) and merges the result into base.
// Keys present in both are resolved according to policy.
func CloneInto(base, src map[string]string, opts Options, policy MergePolicy) MergeResult {
	cloned := Clone(src, opts)

	merged := make(map[string]string, len(base)+len(cloned.Env))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []string
	for k, v := range cloned.Env {
		if existing, ok := merged[k]; ok && existing != v {
			conflicts = append(conflicts, k)
			if policy == PolicyKeepBase {
				continue
			}
		}
		merged[k] = v
	}

	return MergeResult{
		Env:       merged,
		Conflicts: conflicts,
	}
}
