package validator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/validator"
)

func TestValidate_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{}
	rules := []validator.Rule{{Key: "DATABASE_URL", Required: true}}

	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected key DATABASE_URL, got %s", violations[0].Key)
	}
}

func TestValidate_RequiredKeyEmpty(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "   "}
	rules := []validator.Rule{{Key: "DATABASE_URL", Required: true}}

	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_RequiredKeyPresent(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://localhost/db"}
	rules := []validator.Rule{{Key: "DATABASE_URL", Required: true}}

	violations := validator.Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}

	violations := validator.Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "not-a-port"}
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}

	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_OptionalMissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	rules := []validator.Rule{{Key: "OPTIONAL_KEY", Required: false, Pattern: `^\d+$`}}

	violations := validator.Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for missing optional key, got %v", violations)
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "abc",
	}
	rules := []validator.Rule{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
		{Key: "SECRET_KEY", Required: true},
	}

	violations := validator.Validate(env, rules)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
}
