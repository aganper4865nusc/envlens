// Package scoper provides environment variable scoping — filtering and
// projecting an env map to a named scope (e.g. "prod", "staging") based
// on prefix conventions or explicit scope rules.
package scoper

import (
	"sort"
	"strings"
)

// Rule defines how a scope is matched and what transformation to apply.
type Rule struct {
	// Scope is the logical name, e.g. "prod" or "staging".
	Scope string
	// Prefix is the key prefix that identifies this scope, e.g. "PROD_".
	Prefix string
	// StripPrefix removes the prefix from keys in the scoped result.
	StripPrefix bool
}

// Result holds the output of a scoping operation.
type Result struct {
	// Scope is the name of the matched scope.
	Scope string
	// Env contains the keys belonging to the scope.
	Env map[string]string
	// Stripped indicates whether prefixes were removed.
	Stripped bool
	// MatchedKeys lists the original keys that matched.
	MatchedKeys []string
}

// Scope applies the given rule to env and returns a Result.
// Keys not matching the rule's prefix are excluded.
// If rule.StripPrefix is true, the prefix is removed from each key.
func Scope(env map[string]string, rule Rule) Result {
	out := make(map[string]string)
	var matched []string

	for k, v := range env {
		if rule.Prefix == "" || strings.HasPrefix(k, rule.Prefix) {
			matched = append(matched, k)
			newKey := k
			if rule.StripPrefix && rule.Prefix != "" {
				newKey = strings.TrimPrefix(k, rule.Prefix)
			}
			out[newKey] = v
		}
	}

	sort.Strings(matched)

	return Result{
		Scope:       rule.Scope,
		Env:         out,
		Stripped:    rule.StripPrefix && rule.Prefix != "",
		MatchedKeys: matched,
	}
}

// ScopeAll applies multiple rules to env and returns one Result per rule.
// Rules are applied independently; the same key may appear in multiple results.
func ScopeAll(env map[string]string, rules []Rule) []Result {
	results := make([]Result, 0, len(rules))
	for _, r := range rules {
		results = append(results, Scope(env, r))
	}
	return results
}
