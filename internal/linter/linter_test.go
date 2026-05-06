package linter

import (
	"testing"
)

func TestLint_UppercaseKeyWarn(t *testing.T) {
	env := map[string]string{"my_key": "value"}
	issues := Lint(env)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for lowercase key")
	}
	found := false
	for _, i := range issues {
		if i.Key == "my_key" && i.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for lowercase key 'my_key'")
	}
}

func TestLint_ValidKey_NoIssues(t *testing.T) {
	env := map[string]string{"MY_KEY": "value"}
	issues := Lint(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_InvalidKeyChars(t *testing.T) {
	env := map[string]string{"MY-KEY": "value"}
	issues := Lint(env)
	found := false
	for _, i := range issues {
		if i.Key == "MY-KEY" && i.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected error issue for key with invalid characters 'MY-KEY'")
	}
}

func TestLint_TrailingWhitespace(t *testing.T) {
	env := map[string]string{"MY_KEY": "value   "}
	issues := Lint(env)
	found := false
	for _, i := range issues {
		if i.Key == "MY_KEY" && i.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for trailing whitespace in value")
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"bad key": "value",
		"GOOD_KEY": "clean",
		"lower": "val ",
	}
	issues := Lint(env)
	// "bad key" triggers spaces-in-key (error) and valid-key-chars (error) and uppercase (warn)
	// "lower" triggers uppercase (warn) and trailing whitespace (warn)
	if len(issues) < 3 {
		t.Errorf("expected at least 3 issues, got %d", len(issues))
	}
}

func TestLint_EmptyEnv(t *testing.T) {
	issues := Lint(map[string]string{})
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty env, got %d", len(issues))
	}
}
