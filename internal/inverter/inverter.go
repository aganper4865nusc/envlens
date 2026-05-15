// Package inverter provides utilities for inverting environment variable maps,
// swapping keys and values to enable reverse lookups.
package inverter

import "sort"

// Result holds the output of an inversion operation.
type Result struct {
	// Inverted is the new map with values as keys and keys as values.
	Inverted map[string]string
	// Collisions records original keys whose values collided during inversion.
	// The map key is the duplicate value; the slice contains all original keys
	// that shared that value.
	Collisions map[string][]string
}

// Options controls the behaviour of Invert.
type Options struct {
	// OnCollision determines how to handle duplicate values.
	// "first"  – keep the first key seen (alphabetical order).
	// "last"   – keep the last key seen (alphabetical order).
	// "skip"   – omit colliding entries entirely.
	// Default is "first".
	OnCollision string
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{OnCollision: "first"}
}

// Invert swaps keys and values in env. Where multiple keys share the same
// value the collision policy in opts is applied.
func Invert(env map[string]string, opts Options) Result {
	if opts.OnCollision == "" {
		opts.OnCollision = "first"
	}

	// Process keys in a deterministic order.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	inverted := make(map[string]string)
	collisions := make(map[string][]string)

	for _, k := range keys {
		v := env[k]
		if existing, exists := inverted[v]; exists {
			// Record the collision.
			if len(collisions[v]) == 0 {
				collisions[v] = []string{existing}
			}
			collisions[v] = append(collisions[v], k)

			switch opts.OnCollision {
			case "last":
				inverted[v] = k
			case "skip":
				delete(inverted, v)
			// "first" – do nothing; keep the already-stored value.
			}
			continue
		}
		inverted[v] = k
	}

	return Result{
		Inverted:   inverted,
		Collisions: collisions,
	}
}
