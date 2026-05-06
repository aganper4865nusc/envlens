// Package interpolator expands variable references within env file values.
// It supports ${VAR} and $VAR syntax, resolving references from the same
// map or a provided override context.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Result holds the interpolated environment map and any warnings produced
// during expansion (e.g. references to undefined variables).
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Interpolate expands variable references in values of env using ctx as the
// lookup source. If ctx is nil, env itself is used for lookups. References
// that cannot be resolved are left as empty strings and a warning is recorded.
func Interpolate(env map[string]string, ctx map[string]string) Result {
	if ctx == nil {
		ctx = env
	}

	out := make(map[string]string, len(env))
	var warnings []string

	for key, value := range env {
		expanded, w := expand(value, ctx)
		out[key] = expanded
		warnings = append(warnings, w...)
	}

	return Result{Env: out, Warnings: warnings}
}

// expand replaces all variable references in s using the provided lookup map.
func expand(s string, lookup map[string]string) (string, []string) {
	var warnings []string

	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := extractName(match)
		if val, ok := lookup[name]; ok {
			return val
		}
		warnings = append(warnings, fmt.Sprintf("undefined variable: %s", name))
		return ""
	})

	return result, warnings
}

// extractName pulls the variable name out of a $VAR or ${VAR} token.
func extractName(token string) string {
	token = strings.TrimPrefix(token, "$")
	token = strings.TrimPrefix(token, "{")
	token = strings.TrimSuffix(token, "}")
	return token
}
