// Package labeler attaches arbitrary string labels to environment variable keys
// based on configurable rules. Labels are additive — a key may have many labels.
package labeler

import (
	"regexp"
	"sort"
)

// Rule defines a single labeling rule.
type Rule struct {
	// Label is the string to attach when the rule matches.
	Label string
	// Prefix matches keys that start with this string (optional).
	Prefix string
	// Pattern is a regular expression matched against the full key (optional).
	Pattern string
	// Keys is an explicit list of key names to match (optional).
	Keys []string
}

// Result holds the labeled output for the entire env map.
type Result struct {
	// Labels maps each key to its sorted, deduplicated label slice.
	Labels map[string][]string
	// Unlabeled contains keys that matched no rule.
	Unlabeled []string
}

// Label applies rules to env and returns a Result.
func Label(env map[string]string, rules []Rule) (Result, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return Result{}, err
			}
			compiled[i] = re
		}
	}

	labels := make(map[string][]string, len(env))
	explicitSet := make(map[string]map[string]struct{})

	for key := range env {
		explicitSet[key] = make(map[string]struct{})
	}

	for i, rule := range rules {
		keySet := make(map[string]struct{})
		for _, k := range rule.Keys {
			keySet[k] = struct{}{}
		}
		for key := range env {
			matched := false
			if rule.Prefix != "" && len(key) >= len(rule.Prefix) && key[:len(rule.Prefix)] == rule.Prefix {
				matched = true
			}
			if !matched && compiled[i] != nil && compiled[i].MatchString(key) {
				matched = true
			}
			if !matched {
				if _, ok := keySet[key]; ok {
					matched = true
				}
			}
			if matched {
				if _, seen := explicitSet[key][rule.Label]; !seen {
					explicitSet[key][rule.Label] = struct{}{}
					labels[key] = append(labels[key], rule.Label)
				}
			}
		}
	}

	var unlabeled []string
	for key := range env {
		if len(labels[key]) == 0 {
			unlabeled = append(unlabeled, key)
		} else {
			sort.Strings(labels[key])
		}
	}
	sort.Strings(unlabeled)

	return Result{Labels: labels, Unlabeled: unlabeled}, nil
}
