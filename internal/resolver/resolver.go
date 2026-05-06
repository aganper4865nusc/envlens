package resolver

import (
	"fmt"
	"os"
	"strings"
)

// Resolution holds the result of resolving a single env var.
type Resolution struct {
	Key      string
	Value    string
	Source   string // "file", "env", "default"
	Override bool   // true if env overrode the file value
}

// Options controls how resolution behaves.
type Options struct {
	// AllowEnvOverride lets OS environment variables override file values.
	AllowEnvOverride bool
	// Defaults provides fallback values when a key is absent everywhere.
	Defaults map[string]string
}

// Resolve merges a parsed env map with OS environment and defaults,
// returning a slice of Resolution records and any keys that could not
// be resolved.
func Resolve(fileEnv map[string]string, opts Options) ([]Resolution, []string) {
	seen := make(map[string]bool)
	var results []Resolution
	var missing []string

	// Start from file values.
	for k, v := range fileEnv {
		seen[k] = true
		res := Resolution{Key: k, Value: v, Source: "file"}

		if opts.AllowEnvOverride {
			if envVal, ok := os.LookupEnv(k); ok {
				res.Value = envVal
				res.Source = "env"
				res.Override = true
			}
		}
		results = append(results, res)
	}

	// Apply defaults for keys not already seen.
	for k, v := range opts.Defaults {
		if seen[k] {
			continue
		}
		seen[k] = true
		results = append(results, Resolution{
			Key:    k,
			Value:  v,
			Source: "default",
		})
	}

	// Collect keys that exist in OS env but not in file or defaults.
	for _, pair := range os.Environ() {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := parts[0]
		if !seen[k] {
			// Not tracked — not missing, just not in scope.
			_ = k
		}
	}

	// Identify keys in defaults that are still empty after resolution.
	for _, r := range results {
		if r.Value == "" && r.Source == "default" {
			missing = append(missing, fmt.Sprintf("%s (default empty)", r.Key))
		}
	}

	return results, missing
}

// ToMap converts a slice of Resolution into a plain key→value map.
func ToMap(resolutions []Resolution) map[string]string {
	out := make(map[string]string, len(resolutions))
	for _, r := range resolutions {
		out[r.Key] = r.Value
	}
	return out
}
