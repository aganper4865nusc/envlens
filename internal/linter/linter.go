package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a linting issue found in an env file.
type Issue struct {
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// Rule defines a linting rule applied to env key-value pairs.
type Rule struct {
	Name    string
	Check   func(key, value string) *Issue
}

var defaultRules = []Rule{
	{
		Name: "uppercase-key",
		Check: func(key, value string) *Issue {
			if key != strings.ToUpper(key) {
				return &Issue{
					Key:      key,
					Message:  fmt.Sprintf("key %q should be uppercase", key),
					Severity: "warn",
				}
			}
			return nil
		},
	},
	{
		Name: "no-spaces-in-key",
		Check: func(key, value string) *Issue {
			if strings.Contains(key, " ") {
				return &Issue{
					Key:      key,
					Message:  fmt.Sprintf("key %q contains spaces", key),
					Severity: "error",
				}
			}
			return nil
		},
	},
	{
		Name: "valid-key-chars",
		Check: func(key, value string) *Issue {
			validKey := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
			if !validKey.MatchString(key) {
				return &Issue{
					Key:      key,
					Message:  fmt.Sprintf("key %q contains invalid characters", key),
					Severity: "error",
				}
			}
			return nil
		},
	},
	{
		Name: "no-trailing-whitespace",
		Check: func(key, value string) *Issue {
			if value != strings.TrimRight(value, " \t") {
				return &Issue{
					Key:      key,
					Message:  fmt.Sprintf("value for %q has trailing whitespace", key),
					Severity: "warn",
				}
			}
			return nil
		},
	},
}

// Lint runs all default rules against the provided env map and returns any issues found.
func Lint(env map[string]string) []Issue {
	var issues []Issue
	for key, value := range env {
		for _, rule := range defaultRules {
			if issue := rule.Check(key, value); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues
}

// HasErrors returns true if any of the provided issues have severity "error".
func HasErrors(issues []Issue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}
