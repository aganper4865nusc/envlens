package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern for value validation
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s", v.Key, v.Message)
}

// Validate checks the provided env map against a set of rules and returns
// any violations found.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		value, exists := env[rule.Key]

		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.Required && strings.TrimSpace(value) == "" {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is present but empty",
			})
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, value)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !matched {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", value, rule.Pattern),
				})
			}
		}
	}

	return violations
}
