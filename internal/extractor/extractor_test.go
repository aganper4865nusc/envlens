package extractor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/extractor"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestExtract_ExplicitKeys(t *testing.T) {
	src := makeEnv("FOO", "1", "BAR", "2", "BAZ", "3")
	res, err := extractor.Extract(src, extractor.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "1" || res.Env["BAZ"] != "3" {
		t.Errorf("expected FOO=1 and BAZ=3, got %v", res.Env)
	}
	if _, ok := res.Env["BAR"]; ok {
		t.Error("BAR should not be in result")
	}
}

func TestExtract_MissingKeyTracked(t *testing.T) {
	src := makeEnv("FOO", "bar")
	res, err := extractor.Extract(src, extractor.Options{Keys: []string{"FOO", "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING" {
		t.Errorf("expected [MISSING], got %v", res.Missing)
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	src := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	res, err := extractor.Extract(src, extractor.Options{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if _, ok := res.Env["APP_NAME"]; ok {
		t.Error("APP_NAME should not be extracted")
	}
}

func TestExtract_ByPattern(t *testing.T) {
	src := makeEnv("SECRET_KEY", "abc", "API_SECRET", "xyz", "HOST", "localhost")
	res, err := extractor.Extract(src, extractor.Options{Pattern: "SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(res.Env), res.Env)
	}
}

func TestExtract_InvalidPattern_ReturnsError(t *testing.T) {
	src := makeEnv("FOO", "bar")
	_, err := extractor.Extract(src, extractor.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestExtract_NoOptions_ReturnsError(t *testing.T) {
	src := makeEnv("FOO", "bar")
	_, err := extractor.Extract(src, extractor.Options{})
	if err == nil {
		t.Error("expected error when no selectors provided")
	}
}

func TestExtract_ExtractedListIsSorted(t *testing.T) {
	src := makeEnv("Z_KEY", "1", "A_KEY", "2", "M_KEY", "3")
	res, err := extractor.Extract(src, extractor.Options{Prefixes: []string{""}}) // empty prefix matches all
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Extracted); i++ {
		if res.Extracted[i-1] > res.Extracted[i] {
			t.Errorf("extracted list not sorted: %v", res.Extracted)
			break
		}
	}
}
