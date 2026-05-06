// Package trimmer provides utilities for cleaning up environment variable maps
// by removing keys or values that match configurable criteria such as empty values,
// blank keys, or user-supplied key lists.
package trimmer

// Options controls what the Trimmer removes.
type Options struct {
	// RemoveEmpty removes entries whose value is an empty string.
	RemoveEmpty bool
	// RemoveBlankKeys removes entries whose key is empty or whitespace-only.
	RemoveBlankKeys bool
	// RemoveKeys is an explicit list of keys to drop regardless of value.
	RemoveKeys []string
}

// Result holds the output of a Trim operation.
type Result struct {
	// Env is the cleaned environment map.
	Env map[string]string
	// Removed is the list of keys that were dropped, in the order they were encountered.
	Removed []string
}

// Trim applies the given Options to env and returns a Result.
// The original map is never modified.
func Trim(env map[string]string, opts Options) Result {
	removeSet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeSet[k] = struct{}{}
	}

	out := make(map[string]string, len(env))
	var removed []string

	for k, v := range env {
		if opts.RemoveBlankKeys && isBlank(k) {
			removed = append(removed, k)
			continue
		}
		if opts.RemoveEmpty && v == "" {
			removed = append(removed, k)
			continue
		}
		if _, drop := removeSet[k]; drop {
			removed = append(removed, k)
			continue
		}
		out[k] = v
	}

	return Result{Env: out, Removed: removed}
}

func isBlank(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' {
			return false
		}
	}
	return true
}
