// Package extractor provides utilities for extracting a subset of environment
// variables by key list, prefix set, or regex pattern into a new map.
package extractor

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Options controls how keys are selected during extraction.
type Options struct {
	// Keys is an explicit list of keys to extract.
	Keys []string
	// Prefixes extracts all keys that start with any of the given prefixes.
	Prefixes []string
	// Pattern extracts all keys matching the given regular expression.
	Pattern string
}

// Result holds the extracted environment map and metadata.
type Result struct {
	Env      map[string]string
	// Extracted lists the keys that were successfully pulled out.
	Extracted []string
	// Missing lists explicitly requested keys that were absent in the source.
	Missing []string
}

// Extract copies matching key-value pairs from src into a new map according to
// the provided Options. At least one selector (Keys, Prefixes, or Pattern) must
// be set; otherwise an error is returned.
func Extract(src map[string]string, opts Options) (Result, error) {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 && opts.Pattern == "" {
		return Result{}, fmt.Errorf("extractor: at least one of Keys, Prefixes, or Pattern must be specified")
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("extractor: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	out := make(map[string]string)
	selected := make(map[string]bool)

	// Prefix-based selection.
	for k, v := range src {
		for _, p := range opts.Prefixes {
			if strings.HasPrefix(k, p) {
				out[k] = v
				selected[k] = true
				break
			}
		}
	}

	// Pattern-based selection.
	if re != nil {
		for k, v := range src {
			if re.MatchString(k) {
				out[k] = v
				selected[k] = true
			}
		}
	}

	// Explicit key selection; track missing ones.
	var missing []string
	for _, k := range opts.Keys {
		if v, ok := src[k]; ok {
			out[k] = v
			selected[k] = true
		} else {
			missing = append(missing, k)
		}
	}

	extracted := make([]string, 0, len(selected))
	for k := range selected {
		extracted = append(extracted, k)
	}
	sort.Strings(extracted)
	sort.Strings(missing)

	return Result{Env: out, Extracted: extracted, Missing: missing}, nil
}
