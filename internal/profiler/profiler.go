// Package profiler analyzes environment variable sets and produces
// a summary profile: key count, empty values, sensitive key ratio, etc.
package profiler

import (
	"sort"
	"strings"
)

// Profile holds statistical metadata about an env map.
type Profile struct {
	Source        string
	TotalKeys     int
	EmptyValues   int
	SensitiveKeys []string
	UppercaseKeys int
	LowercaseKeys int
	MixedKeys     int
	UniqueValues  int
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

// Analyze builds a Profile from the provided env map and source label.
func Analyze(source string, env map[string]string) Profile {
	p := Profile{Source: source}

	valueSeen := make(map[string]struct{})

	for k, v := range env {
		p.TotalKeys++

		if v == "" {
			p.EmptyValues++
		} else {
			valueSeen[v] = struct{}{}
		}

		if isSensitive(k) {
			p.SensitiveKeys = append(p.SensitiveKeys, k)
		}

		switch keyCase(k) {
		case "upper":
			p.UppercaseKeys++
		case "lower":
			p.LowercaseKeys++
		default:
			p.MixedKeys++
		}
	}

	p.UniqueValues = len(valueSeen)
	sort.Strings(p.SensitiveKeys)
	return p
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pat := range sensitivePatterns {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}

func keyCase(key string) string {
	if key == strings.ToUpper(key) {
		return "upper"
	}
	if key == strings.ToLower(key) {
		return "lower"
	}
	return "mixed"
}
