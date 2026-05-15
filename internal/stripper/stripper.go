// Package stripper removes keys from an env map based on prefix, pattern, or explicit list.
package stripper

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls which keys are stripped from the environment map.
type Options struct {
	// Prefixes removes any key that starts with one of these prefixes.
	Prefixes []string
	// Keys removes these exact keys.
	Keys []string
	// Pattern removes keys matching this regular expression.
	Pattern string
}

// Result holds the output of a Strip operation.
type Result struct {
	// Env is the new map with matched keys removed.
	Env map[string]string
	// Stripped is the list of keys that were removed, in sorted order.
	Stripped []string
}

// Strip removes keys from env according to opts and returns a Result.
// The original map is never modified.
func Strip(env map[string]string, opts Options) (Result, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("stripper: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	out := make(map[string]string, len(env))
	var stripped []string

	for k, v := range env {
		if shouldStrip(k, opts.Prefixes, opts.Keys, re) {
			stripped = append(stripped, k)
			continue
		}
		out[k] = v
	}

	sortStrings(stripped)
	return Result{Env: out, Stripped: stripped}, nil
}

func shouldStrip(key string, prefixes, keys []string, re *regexp.Regexp) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	if re != nil && re.MatchString(key) {
		return true
	}
	return false
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
