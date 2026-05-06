package redactor

import (
	"sort"
	"testing"
)

func TestRedact_SensitiveKeyIsRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "envlens",
	}
	result := Redact(env)

	if result.Env["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result.Env["DB_PASSWORD"])
	}
	if result.Env["APP_NAME"] != "envlens" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result.Env["APP_NAME"])
	}
}

func TestRedact_EmptySensitiveValueNotRedacted(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	result := Redact(env)

	if result.Env["API_KEY"] != "" {
		t.Errorf("expected empty API_KEY to remain empty, got %q", result.Env["API_KEY"])
	}
	if len(result.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", result.RedactedKeys)
	}
}

func TestRedact_MultiplePatterns(t *testing.T) {
	env := map[string]string{
		"AUTH_TOKEN":      "tok_abc123",
		"PRIVATE_KEY":     "-----BEGIN RSA",
		"DATABASE_URL":    "postgres://user:pass@host/db",
		"LOG_LEVEL":       "debug",
		"SECRET_KEY":      "my_secret",
	}
	result := Redact(env)

	expectedRedacted := []string{"AUTH_TOKEN", "PRIVATE_KEY", "DATABASE_URL", "SECRET_KEY"}
	for _, key := range expectedRedacted {
		if result.Env[key] != "[REDACTED]" {
			t.Errorf("expected %s to be redacted, got %q", key, result.Env[key])
		}
	}
	if result.Env["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL to be unchanged")
	}
}

func TestRedact_RedactedKeysTracked(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "pass",
		"API_TOKEN":   "tok",
		"HOST":        "localhost",
	}
	result := Redact(env)

	sort.Strings(result.RedactedKeys)
	if len(result.RedactedKeys) != 2 {
		t.Fatalf("expected 2 redacted keys, got %d: %v", len(result.RedactedKeys), result.RedactedKeys)
	}
	if result.RedactedKeys[0] != "API_TOKEN" || result.RedactedKeys[1] != "DB_PASSWORD" {
		t.Errorf("unexpected redacted keys: %v", result.RedactedKeys)
	}
}

func TestRedact_OriginalMapUnmodified(t *testing.T) {
	env := map[string]string{
		"SECRET": "original_value",
	}
	Redact(env)

	if env["SECRET"] != "original_value" {
		t.Errorf("original map should not be modified")
	}
}
