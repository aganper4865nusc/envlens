// Package filter provides functionality to filter environment variable maps
// by key patterns, prefixes, or custom predicates.
package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Prefix retains only keys that start with the given prefix.
	Prefix string

	// Pattern retains only keys matching the given regular expression.
	Pattern string

	// Invert reverses the filter, retaining keys that do NOT match.
	Invert bool
}

// Result holds the output of a filter operation.
type Result struct {
	Matched  map[string]string
	Excluded map[string]string
}

// Filter applies the given Options to env and returns a Result containing
// matched and excluded keys.
func Filter(env map[string]string, opts Options) (Result, error) {
	matched := make(map[string]string)
	excluded := make(map[string]string)

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, err
		}
	}

	for k, v := range env {
		keep := matches(k, opts.Prefix, re)
		if opts.Invert {
			keep = !keep
		}
		if keep {
			matched[k] = v
		} else {
			excluded[k] = v
		}
	}

	return Result{Matched: matched, Excluded: excluded}, nil
}

func matches(key, prefix string, re *regexp.Regexp) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
