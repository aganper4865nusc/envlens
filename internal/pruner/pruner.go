// Package pruner removes environment variable entries that match
// specified criteria such as key prefix, pattern, or explicit list,
// and optionally keeps only entries that have non-empty values.
package pruner

import (
	"fmt"
	"regexp"
	"sort"
)

// Options controls which keys are pruned from the environment map.
type Options struct {
	// Prefixes removes any key starting with one of these strings.
	Prefixes []string
	// Keys removes these exact keys.
	Keys []string
	// Patterns removes keys matching any of these regular expressions.
	Patterns []string
	// DropEmpty removes keys whose value is the empty string.
	DropEmpty bool
}

// Result holds the outcome of a Prune operation.
type Result struct {
	// Pruned is the new map with matching entries removed.
	Pruned map[string]string
	// Removed lists the keys that were removed, in sorted order.
	Removed []string
}

// Prune returns a copy of env with entries removed according to opts.
// The original map is never modified.
func Prune(env map[string]string, opts Options) (Result, error) {
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return Result{}, fmt.Errorf("pruner: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	out := make(map[string]string, len(env))
	var removed []string

	for k, v := range env {
		if shouldPrune(k, v, opts.Prefixes, explicit, compiled, opts.DropEmpty) {
			removed = append(removed, k)
		} else {
			out[k] = v
		}
	}

	sort.Strings(removed)
	return Result{Pruned: out, Removed: removed}, nil
}

func shouldPrune(key, value string, prefixes []string, explicit map[string]struct{}, patterns []*regexp.Regexp, dropEmpty bool) bool {
	if dropEmpty && value == "" {
		return true
	}
	if _, ok := explicit[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
