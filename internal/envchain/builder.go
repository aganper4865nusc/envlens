package envchain

import (
	"strings"
)

// CommonSteps provides factory functions for frequently used chain steps
// built on top of simple inline transformations.

// StripPrefixStep returns a Step that removes a key prefix from all matching keys.
func StripPrefixStep(prefix string) Step {
	return Step{
		Name: "strip-prefix:" + prefix,
		Fn: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				if strings.HasPrefix(k, prefix) {
					out[strings.TrimPrefix(k, prefix)] = v
				} else {
					out[k] = v
				}
			}
			return out, nil
		},
	}
}

// AddPrefixStep returns a Step that prepends a prefix to every key.
func AddPrefixStep(prefix string) Step {
	return Step{
		Name: "add-prefix:" + prefix,
		Fn: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				out[prefix+k] = v
			}
			return out, nil
		},
	}
}

// DropEmptyStep returns a Step that removes keys with empty values.
func DropEmptyStep() Step {
	return Step{
		Name: "drop-empty",
		Fn: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				if strings.TrimSpace(v) != "" {
					out[k] = v
				}
			}
			return out, nil
		},
	}
}

// UppercaseKeysStep returns a Step that converts all keys to uppercase.
func UppercaseKeysStep() Step {
	return Step{
		Name: "uppercase-keys",
		Fn: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				out[strings.ToUpper(k)] = v
			}
			return out, nil
		},
	}
}
