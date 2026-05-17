// Package mapper provides utilities for transforming environment variable maps
// using key-value mapping rules (e.g. key aliasing, value substitution).
package mapper

import "sort"

// KeyRule defines a mapping from one key name to another.
type KeyRule struct {
	From string
	To   string
}

// ValueRule defines a substitution for a specific key's value.
type ValueRule struct {
	Key  string
	From string
	To   string
}

// Options controls the behaviour of Map.
type Options struct {
	// KeyRules renames keys according to the provided rules.
	KeyRules []KeyRule
	// ValueRules substitutes values for specific keys.
	ValueRules []ValueRule
	// DropUnmapped drops all keys not referenced by any KeyRule.
	DropUnmapped bool
}

// Result holds the output of a Map operation.
type Result struct {
	Env         map[string]string
	RenamedKeys []string
	Substituted []string
}

// Map applies key and value mapping rules to env and returns a Result.
func Map(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))

	// Build lookup tables.
	keyMap := make(map[string]string, len(opts.KeyRules))
	for _, r := range opts.KeyRules {
		keyMap[r.From] = r.To
	}

	valMap := make(map[string]map[string]string)
	for _, r := range opts.ValueRules {
		if valMap[r.Key] == nil {
		valMap[r.Key] = make(map[string]string)
		}
		valMap[r.Key][r.From] = r.To
	}

	var renamed, substituted []string

	for k, v := range env {
		destKey := k
		if mapped, ok := keyMap[k]; ok {
			destKey = mapped
			renamed = append(renamed, k)
		} else if opts.DropUnmapped && len(opts.KeyRules) > 0 {
			continue
		}

		// Apply value substitution using the original key name.
		if subs, ok := valMap[k]; ok {
			if replacement, ok2 := subs[v]; ok2 {
				v = replacement
				substituted = append(substituted, destKey)
			}
		}

		out[destKey] = v
	}

	sort.Strings(renamed)
	sort.Strings(substituted)

	return Result{
		Env:         out,
		RenamedKeys: renamed,
		Substituted: substituted,
	}
}
