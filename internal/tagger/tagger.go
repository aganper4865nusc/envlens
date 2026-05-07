// Package tagger assigns metadata tags to environment variable keys
// based on configurable rules such as prefix matching, pattern matching,
// or explicit key lists.
package tagger

import (
	"regexp"
	"strings"
)

// Rule defines a single tagging rule.
type Rule struct {
	// Tag is the label to apply when this rule matches.
	Tag string
	// Prefixes is a list of key prefixes that trigger this rule.
	Prefixes []string
	// Pattern is an optional regex applied to the key name.
	Pattern string
	// Keys is an explicit list of key names to tag.
	Keys []string
}

// Result holds the tagging output for a single key.
type Result struct {
	Key  string
	Tags []string
}

// Tag applies the given rules to the provided env map and returns
// a slice of Result, one per key, in sorted order.
func Tag(env map[string]string, rules []Rule) []Result {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		if r.Pattern != "" {
			compiled[i] = regexp.MustCompile(r.Pattern)
		}
	}

	keySet := make([]string, 0, len(env))
	for k := range env {
		keySet = append(keySet, k)
	}
	sortStrings(keySet)

	results := make([]Result, 0, len(keySet))
	for _, key := range keySet {
		tags := collectTags(key, rules, compiled)
		results = append(results, Result{Key: key, Tags: tags})
	}
	return results
}

func collectTags(key string, rules []Rule, compiled []*regexp.Regexp) []string {
	seen := map[string]bool{}
	var tags []string
	for i, r := range rules {
		if matchesRule(key, r, compiled[i]) {
			if !seen[r.Tag] {
				seen[r.Tag] = true
				tags = append(tags, r.Tag)
			}
		}
	}
	return tags
}

func matchesRule(key string, r Rule, re *regexp.Regexp) bool {
	for _, prefix := range r.Prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	for _, k := range r.Keys {
		if k == key {
			return true
		}
	}
	if re != nil && re.MatchString(key) {
		return true
	}
	return false
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
