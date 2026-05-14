// Package summarizer produces a human-readable summary of an env map,
// reporting counts, value statistics, and category breakdowns.
package summarizer

import (
	"fmt"
	"sort"
	"strings"
)

// Summary holds aggregate statistics for an environment map.
type Summary struct {
	TotalKeys     int
	EmptyValues   int
	NonEmptyValues int
	UniqueValues  int
	Categories    map[string]int // category label -> key count
	LongestKey    string
	LongestValue  string
	Lines         []string // formatted lines for display
}

var categoryPatterns = []struct {
	label   string
	prefixes []string
}{
	{"database", []string{"DB_", "DATABASE_", "POSTGRES_", "MYSQL_", "MONGO_"}},
	{"auth", []string{"AUTH_", "JWT_", "OAUTH_", "SECRET_", "TOKEN_", "API_KEY"}},
	{"network", []string{"HOST", "PORT", "URL", "ADDR", "ENDPOINT"}},
	{"observability", []string{"LOG_", "TRACE_", "METRIC_", "SENTRY_", "DATADOG_"}},
}

// Summarize computes a Summary for the given env map.
func Summarize(env map[string]string) Summary {
	s := Summary{
		Categories: make(map[string]int),
	}

	valueSeen := make(map[string]struct{})

	for k, v := range env {
		s.TotalKeys++
		if v == "" {
			s.EmptyValues++
		} else {
			s.NonEmptyValues++
		}

		if _, seen := valueSeen[v]; !seen && v != "" {
			valueSeen[v] = struct{}{}
			s.UniqueValues++
		}

		if len(k) > len(s.LongestKey) {
			s.LongestKey = k
		}
		if len(v) > len(s.LongestValue) {
			s.LongestValue = v
		}

		upper := strings.ToUpper(k)
		for _, cat := range categoryPatterns {
			for _, pfx := range cat.prefixes {
				if strings.Contains(upper, pfx) {
					s.Categories[cat.label]++
					break
				}
			}
		}
	}

	s.Lines = format(s)
	return s
}

func format(s Summary) []string {
	lines := []string{
		fmt.Sprintf("Total keys    : %d", s.TotalKeys),
		fmt.Sprintf("Non-empty     : %d", s.NonEmptyValues),
		fmt.Sprintf("Empty values  : %d", s.EmptyValues),
		fmt.Sprintf("Unique values : %d", s.UniqueValues),
		fmt.Sprintf("Longest key   : %s (%d chars)", s.LongestKey, len(s.LongestKey)),
		fmt.Sprintf("Longest value : %d chars", len(s.LongestValue)),
	}

	if len(s.Categories) > 0 {
		lines = append(lines, "Categories:")
		keys := make([]string, 0, len(s.Categories))
		for k := range s.Categories {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			lines = append(lines, fmt.Sprintf("  %-16s %d", k, s.Categories[k]))
		}
	}

	return lines
}
