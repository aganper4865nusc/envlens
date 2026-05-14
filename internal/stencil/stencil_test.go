package stencil_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/stencil"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestGenerate_SensitiveKeyGetsMaskedPlaceholder(t *testing.T) {
	env := makeEnv("DB_PASSWORD", "hunter2")
	entries := stencil.Generate(env, stencil.DefaultOptions())
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Sensitive {
		t.Error("expected DB_PASSWORD to be marked sensitive")
	}
	if entries[0].Placeholder == "hunter2" {
		t.Error("sensitive value should not appear in placeholder")
	}
}

func TestGenerate_NonSensitiveKeyPlaceholder(t *testing.T) {
	env := makeEnv("APP_PORT", "8080")
	entries := stencil.Generate(env, stencil.DefaultOptions())
	if entries[0].Sensitive {
		t.Error("APP_PORT should not be sensitive")
	}
	if entries[0].Placeholder == "8080" {
		t.Error("value should be replaced by placeholder by default")
	}
}

func TestGenerate_PreserveDefaults(t *testing.T) {
	opts := stencil.DefaultOptions()
	opts.PreserveDefaults = true
	env := makeEnv("LOG_LEVEL", "info")
	entries := stencil.Generate(env, opts)
	if entries[0].Placeholder != "info" {
		t.Errorf("expected preserved value 'info', got %q", entries[0].Placeholder)
	}
}

func TestGenerate_EmptyValuePreserved(t *testing.T) {
	env := makeEnv("OPTIONAL_FLAG", "")
	entries := stencil.Generate(env, stencil.DefaultOptions())
	if entries[0].Placeholder != "" {
		t.Errorf("empty value should stay empty, got %q", entries[0].Placeholder)
	}
}

func TestGenerate_SortedOutput(t *testing.T) {
	env := makeEnv("Z_KEY", "z", "A_KEY", "a", "M_KEY", "m")
	entries := stencil.Generate(env, stencil.DefaultOptions())
	if entries[0].Key != "A_KEY" || entries[1].Key != "M_KEY" || entries[2].Key != "Z_KEY" {
		t.Error("entries should be sorted alphabetically")
	}
}

func TestRender_ContainsKeysAndPlaceholders(t *testing.T) {
	env := makeEnv("API_TOKEN", "secret", "APP_ENV", "production")
	opts := stencil.DefaultOptions()
	entries := stencil.Generate(env, opts)
	out := stencil.Render(entries, opts.CommentPrefix)
	if !strings.Contains(out, "API_TOKEN=") {
		t.Error("output should contain API_TOKEN=")
	}
	if strings.Contains(out, "secret") {
		t.Error("sensitive value 'secret' should not appear in rendered output")
	}
}

func TestRender_SensitiveKeyHasComment(t *testing.T) {
	env := makeEnv("DB_SECRET", "abc123")
	opts := stencil.DefaultOptions()
	entries := stencil.Generate(env, opts)
	out := stencil.Render(entries, opts.CommentPrefix)
	if !strings.Contains(out, "Sensitive") {
		t.Error("sensitive key should produce a comment in rendered output")
	}
}
