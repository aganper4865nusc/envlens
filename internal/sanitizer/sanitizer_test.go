package sanitizer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/sanitizer"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestSanitize_TrimWhitespace(t *testing.T) {
	env := makeEnv("KEY", "  hello  ")
	opts := sanitizer.DefaultOptions()
	res := sanitizer.Sanitize(env, opts)
	if got := res.Env["KEY"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
	if len(res.Changed) != 1 || res.Changed[0] != "KEY" {
		t.Errorf("expected KEY in Changed, got %v", res.Changed)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	env := makeEnv("KEY", "val\x00ue\x1F")
	opts := sanitizer.DefaultOptions()
	res := sanitizer.Sanitize(env, opts)
	if got := res.Env["KEY"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestSanitize_CollapseWhitespace(t *testing.T) {
	env := makeEnv("KEY", "hello   world")
	opts := sanitizer.Options{CollapseWhitespace: true, TrimWhitespace: false}
	res := sanitizer.Sanitize(env, opts)
	if got := res.Env["KEY"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestSanitize_NoChanges(t *testing.T) {
	env := makeEnv("KEY", "clean_value")
	opts := sanitizer.DefaultOptions()
	res := sanitizer.Sanitize(env, opts)
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
}

func TestSanitize_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("KEY", "  dirty  ")
	opts := sanitizer.DefaultOptions()
	sanitizer.Sanitize(env, opts)
	if env["KEY"] != "  dirty  " {
		t.Error("original map was modified")
	}
}

func TestSanitize_ChangedKeysSorted(t *testing.T) {
	env := makeEnv("ZEBRA", "  z  ", "ALPHA", "  a  ", "MIDDLE", "  m  ")
	opts := sanitizer.DefaultOptions()
	res := sanitizer.Sanitize(env, opts)
	if len(res.Changed) != 3 {
		t.Fatalf("expected 3 changed keys, got %d", len(res.Changed))
	}
	if res.Changed[0] != "ALPHA" || res.Changed[1] != "MIDDLE" || res.Changed[2] != "ZEBRA" {
		t.Errorf("expected sorted order, got %v", res.Changed)
	}
}

func TestSanitize_RemoveNonPrintable(t *testing.T) {
	env := makeEnv("KEY", "val\u200Bue") // zero-width space
	opts := sanitizer.Options{RemoveNonPrintable: true}
	res := sanitizer.Sanitize(env, opts)
	if got := res.Env["KEY"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}
