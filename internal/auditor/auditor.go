package auditor

import (
	"fmt"
	"sort"
	"strings"
)

// AuditResult holds the findings for a single key in an env file.
type AuditResult struct {
	Key      string
	Severity string // "warn" or "info"
	Message  string
}

// AuditOptions configures which checks the auditor performs.
type AuditOptions struct {
	// FlagPlainSecrets checks for keys whose names suggest secrets but have plain-text values.
	FlagPlainSecrets bool
	// FlagEmptyValues reports keys that have empty values.
	FlagEmptyValues bool
	// FlagDuplicateKeys reports duplicate key definitions (parser keeps last; we surface it).
	FlagDuplicateKeys bool
}

var secretKeyHints = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY", "PRIVATE_KEY",
}

// Audit inspects an env map and returns a list of audit findings.
// rawLines is the original slice of non-comment, non-blank lines from the file
// (used to detect duplicate keys). Pass nil to skip duplicate detection.
func Audit(env map[string]string, rawLines []string, opts AuditOptions) []AuditResult {
	var results []AuditResult

	if opts.FlagDuplicateKeys && rawLines != nil {
		seen := map[string]int{}
		for _, line := range rawLines {
			if idx := strings.Index(line, "="); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				seen[key]++
			}
		}
		for key, count := range seen {
			if count > 1 {
				results = append(results, AuditResult{
					Key:      key,
					Severity: "warn",
					Message:  fmt.Sprintf("key defined %d times; only the last value is used", count),
				})
			}
		}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := env[key]

		if opts.FlagEmptyValues && val == "" {
			results = append(results, AuditResult{
				Key:      key,
				Severity: "warn",
				Message:  "value is empty",
			})
		}

		if opts.FlagPlainSecrets {
			upper := strings.ToUpper(key)
			for _, hint := range secretKeyHints {
				if strings.Contains(upper, hint) {
					if val != "" && !looksEncoded(val) {
						results = append(results, AuditResult{
							Key:      key,
							Severity: "warn",
							Message:  "key name suggests a secret but value appears to be plain text",
						})
					}
					break
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// looksEncoded is a heuristic: values longer than 16 chars with mixed case
// and digits are assumed to be hashed/encoded secrets.
func looksEncoded(val string) bool {
	if len(val) < 16 {
		return false
	}
	hasUpper, hasLower, hasDigit := false, false, false
	for _, c := range val {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}
