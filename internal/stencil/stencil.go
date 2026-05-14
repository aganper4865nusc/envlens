// Package stencil generates a template .env file from an existing env map,
// replacing values with typed placeholders for documentation or onboarding.
package stencil

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how placeholders are generated.
type Options struct {
	// PreserveDefaults keeps non-sensitive values as-is instead of replacing them.
	PreserveDefaults bool
	// SensitivePatterns is a list of substrings that mark a key as sensitive.
	SensitivePatterns []string
	// CommentPrefix is prepended to each placeholder comment line.
	CommentPrefix string
}

// DefaultOptions returns sensible defaults for stencil generation.
func DefaultOptions() Options {
	return Options{
		PreserveDefaults:  false,
		SensitivePatterns: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE"},
		CommentPrefix:     "# ",
	}
}

// Entry represents a single line in the generated stencil.
type Entry struct {
	Key         string
	Placeholder string
	Sensitive   bool
}

// Generate produces a slice of stencil entries from the provided env map.
func Generate(env map[string]string, opts Options) []Entry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		sensitive := isSensitiveKey(k, opts.SensitivePatterns)
		placeholder := buildPlaceholder(k, v, sensitive, opts)
		entries = append(entries, Entry{
			Key:         k,
			Placeholder: placeholder,
			Sensitive:   sensitive,
		})
	}
	return entries
}

// Render formats entries into a .env-style string.
func Render(entries []Entry, commentPrefix string) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Sensitive {
			fmt.Fprintf(&sb, "%sSensitive — do not commit real value\n", commentPrefix)
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Placeholder)
	}
	return sb.String()
}

func isSensitiveKey(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func buildPlaceholder(key, value string, sensitive bool, opts Options) string {
	if sensitive {
		return fmt.Sprintf("<your_%s>", strings.ToLower(key))
	}
	if opts.PreserveDefaults && value != "" {
		return value
	}
	if value == "" {
		return ""
	}
	return fmt.Sprintf("<your_%s>", strings.ToLower(key))
}
