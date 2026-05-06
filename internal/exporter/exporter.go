// Package exporter converts parsed env maps into various output formats
// such as shell export statements, Docker --env-file format, and JSON.
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
)

// Export writes the env map to w in the specified format.
// Keys are written in sorted order for deterministic output.
func Export(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)

	switch format {
	case FormatShell:
		return exportShell(w, env, keys)
	case FormatDocker:
		return exportDocker(w, env, keys)
	case FormatJSON:
		return exportJSON(w, env, keys)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

// exportShell writes lines like: export KEY="value"
func exportShell(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		v := shellEscape(env[k])
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

// exportDocker writes lines like: KEY=value (no quotes, Docker env-file format)
func exportDocker(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

// exportJSON writes a JSON object of key/value pairs.
func exportJSON(w io.Writer, env map[string]string, keys []string) error {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = env[k]
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

// shellEscape wraps value in double quotes and escapes internal double quotes.
func shellEscape(v string) string {
	escaped := strings.ReplaceAll(v, `"`, `\"`)
	return `"` + escaped + `"`
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
