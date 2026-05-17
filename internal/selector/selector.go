// Package selector provides utilities for selecting a subset of environment
// variables based on explicit keys, prefixes, or regex patterns.
package selector

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Options controls how selection is performed.
type Options struct {
	// Keys is a list of exact key names to select.
	Keys []string
	// Prefixes selects all keys that start with any of the given prefixes.
	Prefixes []string
	// Patterns selects all keys matching any of the given regular expressions.
	Patterns []string
	// Invert returns keys that do NOT match the criteria.
	Invert bool
}

// Result holds the outcome of a selection operation.
type Result struct {
	// Env is the selected (or excluded) subset of the environment.
	Env map[string]string
	// Selected contains the keys that matched the criteria.
	Selected []string
	// Skipped contains the keys that did not match.
	Skipped []string
}

// Select applies the given Options to env and returns a Result.
func Select(env map[string]string, opts Options) (Result, error) {
	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return Result{}, fmt.Errorf("selector: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	result := Result{
		Env: make(map[string]string),
	}

	for k, v := range env {
		matched := matchesAny(k, keySet, opts.Prefixes, compiled)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result.Env[k] = v
			result.Selected = append(result.Selected, k)
		} else {
			result.Skipped = append(result.Skipped, k)
		}
	}

	sort.Strings(result.Selected)
	sort.Strings(result.Skipped)
	return result, nil
}

func matchesAny(key string, keySet map[string]struct{}, prefixes []string, patterns []*regexp.Regexp) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
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
