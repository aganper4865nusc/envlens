// Package splitter splits an environment map into named buckets
// based on configurable prefix or pattern rules.
package splitter

import (
	"regexp"
	"sort"
	"strings"
)

// Rule defines how keys are assigned to a named bucket.
type Rule struct {
	// Name is the bucket label.
	Name string
	// Prefix matches keys that start with this string (case-sensitive).
	Prefix string
	// Pattern is a regex alternative to Prefix. If set, Prefix is ignored.
	Pattern string
}

// Result holds the output of a Split operation.
type Result struct {
	// Buckets maps bucket name → key/value pairs assigned to that bucket.
	Buckets map[string]map[string]string
	// Unmatched contains keys that did not match any rule.
	Unmatched map[string]string
}

// Split partitions env according to rules. Each key is placed in the
// first matching bucket. Keys that match no rule go to Unmatched.
func Split(env map[string]string, rules []Rule) (Result, error) {
	result := Result{
		Buckets:   make(map[string]map[string]string),
		Unmatched: make(map[string]string),
	}

	type compiled struct {
		rule Rule
		re   *regexp.Regexp
	}

	var matchers []compiled
	for _, r := range rules {
		result.Buckets[r.Name] = make(map[string]string)
		var re *regexp.Regexp
		if r.Pattern != "" {
			var err error
			re, err = regexp.Compile(r.Pattern)
			if err != nil {
				return Result{}, err
			}
		}
		matchers = append(matchers, compiled{rule: r, re: re})
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		matched := false
		for _, m := range matchers {
			if m.re != nil {
				if m.re.MatchString(k) {
					result.Buckets[m.rule.Name][k] = v
					matched = true
					break
				}
			} else if m.rule.Prefix != "" && strings.HasPrefix(k, m.rule.Prefix) {
				result.Buckets[m.rule.Name][k] = v
				matched = true
				break
			}
		}
		if !matched {
			result.Unmatched[k] = v
		}
	}

	return result, nil
}
