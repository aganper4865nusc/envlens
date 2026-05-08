// Package flattener converts nested prefix-grouped env vars into a flat
// map and vice versa, using a configurable delimiter.
package flattener

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls flattening behaviour.
type Options struct {
	// Delimiter separates prefix segments (default "__").
	Delimiter string
	// MaxDepth limits how many segments are joined (0 = unlimited).
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Delimiter: "__", MaxDepth: 0}
}

// FlatResult holds the flattened map and metadata about what changed.
type FlatResult struct {
	Env       map[string]string
	// Collapsed maps each original key to the key it was stored under.
	Collapsed map[string]string
}

// Flatten takes an env map and collapses keys that share a common prefix
// segment (split by opts.Delimiter) beyond MaxDepth into a single key,
// joining the excess segments back together.
//
// When MaxDepth == 0 the map is returned as-is (no collapsing).
func Flatten(env map[string]string, opts Options) FlatResult {
	if opts.Delimiter == "" {
		opts.Delimiter = "__"
	}

	result := FlatResult{
		Env:       make(map[string]string, len(env)),
		Collapsed: make(map[string]string),
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		newKey := collapse(k, opts.Delimiter, opts.MaxDepth)
		if newKey != k {
			result.Collapsed[k] = newKey
		}
		// Last writer wins on collision.
		result.Env[newKey] = v
	}

	return result
}

// Expand is the inverse of Flatten: given a flat map and a mapping of
// shortKey -> []subKeys it re-expands each entry into multiple keys.
// If expansions is nil the map is returned unchanged.
func Expand(env map[string]string, expansions map[string][]string, delimiter string) map[string]string {
	if delimiter == "" {
		delimiter = "__"
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		subs, ok := expansions[k]
		if !ok || len(subs) == 0 {
			out[k] = v
			continue
		}
		for i, sub := range subs {
			out[sub] = fmt.Sprintf("%s[%d]", v, i)
		}
	}
	return out
}

// collapse joins the first MaxDepth segments and appends the rest as-is.
func collapse(key, delimiter string, maxDepth int) string {
	if maxDepth <= 0 {
		return key
	}
	parts := strings.Split(key, delimiter)
	if len(parts) <= maxDepth {
		return key
	}
	head := strings.Join(parts[:maxDepth], delimiter)
	tail := strings.Join(parts[maxDepth:], delimiter)
	return head + delimiter + tail
}
