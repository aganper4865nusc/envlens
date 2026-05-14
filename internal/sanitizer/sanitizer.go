// Package sanitizer provides utilities for cleaning and normalizing
// environment variable values by removing unsafe or unwanted characters.
package sanitizer

import (
	"regexp"
	"strings"
)

// Options controls the behavior of the sanitizer.
type Options struct {
	// StripControlChars removes ASCII control characters (0x00–0x1F, 0x7F).
	StripControlChars bool
	// TrimWhitespace trims leading and trailing whitespace from values.
	TrimWhitespace bool
	// CollapseWhitespace replaces internal runs of whitespace with a single space.
	CollapseWhitespace bool
	// RemoveNonPrintable removes non-printable Unicode characters.
	RemoveNonPrintable bool
}

// DefaultOptions returns a sensible default sanitizer configuration.
func DefaultOptions() Options {
	return Options{
		StripControlChars:  true,
		TrimWhitespace:     true,
		CollapseWhitespace: false,
		RemoveNonPrintable: false,
	}
}

// Result holds the sanitized environment and a record of which keys were changed.
type Result struct {
	Env     map[string]string
	Changed []string
}

var (
	controlRe      = regexp.MustCompile(`[\x00-\x1F\x7F]`)
	nonPrintableRe = regexp.MustCompile(`[^\x20-\x7E]`)
	collapseRe     = regexp.MustCompile(`\s{2,}`)
)

// Sanitize applies the given options to each value in env, returning a new
// map with cleaned values and a list of keys whose values were modified.
func Sanitize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	var changed []string

	for k, v := range env {
		cleaned := sanitizeValue(v, opts)
		out[k] = cleaned
		if cleaned != v {
			changed = append(changed, k)
		}
	}

	sortStrings(changed)
	return Result{Env: out, Changed: changed}
}

func sanitizeValue(v string, opts Options) string {
	if opts.StripControlChars {
		v = controlRe.ReplaceAllString(v, "")
	}
	if opts.RemoveNonPrintable {
		v = nonPrintableRe.ReplaceAllString(v, "")
	}
	if opts.CollapseWhitespace {
		v = collapseRe.ReplaceAllString(v, " ")
	}
	if opts.TrimWhitespace {
		v = strings.TrimSpace(v)
	}
	return v
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
