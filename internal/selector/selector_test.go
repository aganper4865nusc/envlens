package selector_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/selector"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestSelect_ByExplicitKey(t *testing.T) {
	env := makeEnv("FOO", "1", "BAR", "2", "BAZ", "3")
	res, err := selector.Select(env, selector.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if res.Env["FOO"] != "1" || res.Env["BAZ"] != "3" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestSelect_ByPrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	res, err := selector.Select(env, selector.Options{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if _, ok := res.Env["APP_NAME"]; ok {
		t.Error("APP_NAME should not be selected")
	}
}

func TestSelect_ByPattern(t *testing.T) {
	env := makeEnv("AWS_KEY", "k", "AWS_SECRET", "s", "GCP_TOKEN", "t")
	res, err := selector.Select(env, selector.Options{Patterns: []string{"^AWS_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestSelect_Inverted(t *testing.T) {
	env := makeEnv("SECRET_KEY", "x", "PUBLIC_URL", "http://example.com")
	res, err := selector.Select(env, selector.Options{
		Prefixes: []string{"SECRET_"},
		Invert:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["PUBLIC_URL"]; !ok {
		t.Error("PUBLIC_URL should be selected when inverted")
	}
	if _, ok := res.Env["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should be excluded when inverted")
	}
}

func TestSelect_InvalidPattern_ReturnsError(t *testing.T) {
	env := makeEnv("FOO", "bar")
	_, err := selector.Select(env, selector.Options{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestSelect_SelectedAndSkippedSorted(t *testing.T) {
	env := makeEnv("C", "3", "A", "1", "B", "2")
	res, err := selector.Select(env, selector.Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Selected) != 2 || res.Selected[0] != "A" || res.Selected[1] != "C" {
		t.Errorf("expected sorted selected [A C], got %v", res.Selected)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected skipped [B], got %v", res.Skipped)
	}
}

func TestSelect_EmptyOptions_SelectsNothing(t *testing.T) {
	env := makeEnv("FOO", "1", "BAR", "2")
	res, err := selector.Select(env, selector.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected 0 keys with empty options, got %d", len(res.Env))
	}
}
