// Package walker provides utilities for recursively scanning directories
// and collecting .env files matching configurable patterns.
package walker

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options controls how the directory walk behaves.
type Options struct {
	// Patterns is a list of filename patterns to match (e.g. ".env", "*.env", ".env.*").
	// If empty, defaults to [".env", "*.env", ".env.*"].
	Patterns []string

	// MaxDepth limits recursion depth. 0 means unlimited.
	MaxDepth int

	// SkipHidden skips directories whose names begin with a dot.
	SkipHidden bool
}

// Result holds a discovered file path and its depth relative to the root.
type Result struct {
	Path  string
	Depth int
}

var defaultPatterns = []string{".env", "*.env", ".env.*"}

// Walk traverses root and returns all env files matching opts.
func Walk(root string, opts Options) ([]Result, error) {
	if len(opts.Patterns) == 0 {
		opts.Patterns = defaultPatterns
	}

	var results []Result

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}

		depth := depthOf(rel)

		if d.IsDir() {
			if path == root {
				return nil
			}
			if opts.SkipHidden && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
				return filepath.SkipDir
			}
			return nil
		}

		if matchesAny(d.Name(), opts.Patterns) {
			results = append(results, Result{Path: path, Depth: depth})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})
	return results, nil
}

func matchesAny(name string, patterns []string) bool {
	for _, p := range patterns {
		matched, err := filepath.Match(p, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func depthOf(rel string) int {
	if rel == "." {
		return 0
	}
	return strings.Count(rel, string(os.PathSeparator)) + 1
}
