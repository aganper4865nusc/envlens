package redactor

import (
	"regexp"
	"strings"
)

// sensitivePatterns holds key name patterns that suggest sensitive values.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|api_?key)`),
	regexp.MustCompile(`(?i)(private_?key|priv_key)`),
	regexp.MustCompile(`(?i)(auth|credential|cert)`),
	regexp.MustCompile(`(?i)(dsn|database_url|connection_string)`),
}

const redactedPlaceholder = "[REDACTED]"

// Result holds the redacted environment map and metadata.
type Result struct {
	Env          map[string]string
	RedactedKeys []string
}

// Redact returns a copy of env with sensitive values replaced by a placeholder.
// Keys are identified as sensitive based on their names matching known patterns.
func Redact(env map[string]string) Result {
	redacted := make(map[string]string, len(env))
	var redactedKeys []string

	for k, v := range env {
		if isSensitiveKey(k) && strings.TrimSpace(v) != "" {
			redacted[k] = redactedPlaceholder
			redactedKeys = append(redactedKeys, k)
		} else {
			redacted[k] = v
		}
	}

	return Result{
		Env:          redacted,
		RedactedKeys: redactedKeys,
	}
}

// isSensitiveKey returns true if the key name matches any sensitive pattern.
func isSensitiveKey(key string) bool {
	for _, re := range sensitivePatterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
