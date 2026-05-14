package censor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/censor"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCensor_SensitiveKeyIsReplaced(t *testing.T) {
	env := makeEnv("DB_PASSWORD", "s3cr3t", "APP_NAME", "myapp")
	res := censor.Censor(env, censor.Options{})
	if res.Env["DB_PASSWORD"] != "[CENSORED]" {
		t.Errorf("expected DB_PASSWORD to be censored, got %q", res.Env["DB_PASSWORD"])
	}
	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", res.Env["APP_NAME"])
	}
}

func TestCensor_EmptyValueNotCensored(t *testing.T) {
	env := makeEnv("API_SECRET", "")
	res := censor.Censor(env, censor.Options{})
	if res.Env["API_SECRET"] != "" {
		t.Errorf("expected empty value to remain empty, got %q", res.Env["API_SECRET"])
	}
	if len(res.Censored) != 0 {
		t.Errorf("expected no censored keys for empty value, got %v", res.Censored)
	}
}

func TestCensor_CustomToken(t *testing.T) {
	env := makeEnv("AUTH_TOKEN", "abc123")
	res := censor.Censor(env, censor.Options{Token: "***"})
	if res.Env["AUTH_TOKEN"] != "***" {
		t.Errorf("expected custom token, got %q", res.Env["AUTH_TOKEN"])
	}
}

func TestCensor_ValuePatternMatch(t *testing.T) {
	env := makeEnv("SOME_VAR", "Bearer eyJhbGci")
	res := censor.Censor(env, censor.Options{
		ValuePatterns: []string{`(?i)^bearer\s`},
	})
	if res.Env["SOME_VAR"] != "[CENSORED]" {
		t.Errorf("expected value pattern to trigger censor, got %q", res.Env["SOME_VAR"])
	}
}

func TestCensor_CensoredKeysTracked(t *testing.T) {
	env := makeEnv("DB_SECRET", "x", "DB_PASSWORD", "y", "HOST", "localhost")
	res := censor.Censor(env, censor.Options{})
	if len(res.Censored) != 2 {
		t.Errorf("expected 2 censored keys, got %d: %v", len(res.Censored), res.Censored)
	}
}

func TestCensor_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("API_KEY", "original")
	censor.Censor(env, censor.Options{})
	if env["API_KEY"] != "original" {
		t.Errorf("original map was modified")
	}
}

func TestCensor_CustomKeyPattern(t *testing.T) {
	env := makeEnv("INTERNAL_PASS", "topsecret", "PUBLIC_KEY", "not-sensitive")
	res := censor.Censor(env, censor.Options{
		KeyPatterns: []string{`(?i)_pass$`},
	})
	if res.Env["INTERNAL_PASS"] != "[CENSORED]" {
		t.Errorf("expected INTERNAL_PASS censored")
	}
	if res.Env["PUBLIC_KEY"] != "not-sensitive" {
		t.Errorf("expected PUBLIC_KEY unchanged")
	}
}
