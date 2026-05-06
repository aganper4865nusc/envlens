// Package masker provides utilities for partially masking environment variable
// values for safe display in logs, reports, and terminal output.
package masker

import "strings"

// Style controls how values are masked.
type Style int

const (
	// StyleFull replaces the entire value with asterisks.
	StyleFull Style = iota
	// StylePartial reveals the first and last characters, masking the middle.
	StylePartial
	// StylePrefix reveals only the first few characters.
	StylePrefix
)

// Options configures masking behaviour.
type Options struct {
	Style       Style
	PrefixLen   int // characters to reveal at the start (StylePrefix / StylePartial)
	SuffixLen   int // characters to reveal at the end  (StylePartial only)
	MaskChar    rune
	MinLength   int // values shorter than this are masked in full regardless of style
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Style:     StylePartial,
		PrefixLen: 3,
		SuffixLen: 3,
		MaskChar:  '*',
		MinLength: 8,
	}
}

// Mask applies the masking options to a single value.
func Mask(value string, opts Options) string {
	if value == "" {
		return ""
	}

	maskStr := func(n int) string {
		return strings.Repeat(string(opts.MaskChar), n)
	}

	switch opts.Style {
	case StyleFull:
		return maskStr(len(value))

	case StylePrefix:
		if len(value) <= opts.MinLength || len(value) <= opts.PrefixLen {
			return maskStr(len(value))
		}
		return value[:opts.PrefixLen] + maskStr(len(value)-opts.PrefixLen)

	case StylePartial:
		visible := opts.PrefixLen + opts.SuffixLen
		if len(value) <= opts.MinLength || len(value) <= visible {
			return maskStr(len(value))
		}
		midLen := len(value) - visible
		return value[:opts.PrefixLen] + maskStr(midLen) + value[len(value)-opts.SuffixLen:]
	}

	return maskStr(len(value))
}

// MaskMap returns a new map with the specified keys masked using the given options.
func MaskMap(env map[string]string, keys []string, opts Options) map[string]string {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if _, sensitive := keySet[k]; sensitive {
			result[k] = Mask(v, opts)
		} else {
			result[k] = v
		}
	}
	return result
}
