package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/masker"
)

func TestMask_EmptyValue(t *testing.T) {
	opts := masker.DefaultOptions()
	if got := masker.Mask("", opts); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_FullStyle(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.Style = masker.StyleFull
	got := masker.Mask("supersecret", opts)
	if got != "***********" {
		t.Errorf("expected all asterisks, got %q", got)
	}
}

func TestMask_PartialStyle_LongValue(t *testing.T) {
	opts := masker.DefaultOptions() // PrefixLen=3, SuffixLen=3, MinLength=8
	got := masker.Mask("abcdefghij", opts)
	// expect: abc****hij
	if !strings.HasPrefix(got, "abc") {
		t.Errorf("expected prefix 'abc', got %q", got)
	}
	if !strings.HasSuffix(got, "hij") {
		t.Errorf("expected suffix 'hij', got %q", got)
	}
	if !strings.Contains(got, "****") {
		t.Errorf("expected masked middle, got %q", got)
	}
}

func TestMask_PartialStyle_ShortValue_MaskedFully(t *testing.T) {
	opts := masker.DefaultOptions()
	got := masker.Mask("abc", opts) // shorter than MinLength=8
	if got != "***" {
		t.Errorf("expected fully masked, got %q", got)
	}
}

func TestMask_PrefixStyle(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.Style = masker.StylePrefix
	opts.PrefixLen = 4
	got := masker.Mask("token_abc123xyz", opts)
	if !strings.HasPrefix(got, "toke") {
		t.Errorf("expected prefix 'toke', got %q", got)
	}
	if got[4:] != strings.Repeat("*", len("token_abc123xyz")-4) {
		t.Errorf("expected masked suffix, got %q", got)
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.Style = masker.StyleFull
	opts.MaskChar = '#'
	got := masker.Mask("hello", opts)
	if got != "#####" {
		t.Errorf("expected '#####', got %q", got)
	}
}

func TestMask_OutputLengthMatchesInput(t *testing.T) {
	opts := masker.DefaultOptions()
	values := []string{"a", "ab", "abcdefgh", "supersecrettoken123"}
	for _, v := range values {
		got := masker.Mask(v, opts)
		if len(got) != len(v) {
			t.Errorf("Mask(%q): output length %d does not match input length %d", v, len(got), len(v))
		}
	}
}

func TestMaskMap_OnlySensitiveKeysMasked(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecretpass",
		"API_KEY":     "key_abc123xyz789",
	}
	opts := masker.DefaultOptions()
	result := masker.MaskMap(env, []string{"DB_PASSWORD", "API_KEY"}, opts)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", result["APP_NAME"])
	}
	if result["DB_PASSWORD"] == "supersecretpass" {
		t.Error("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] == "key_abc123xyz789" {
		t.Error("API_KEY should be masked")
	}
}

func TestMaskMap_OriginalUnmodified(t *testing.T) {
	env := map[string]string{"SECRET": "plaintext"}
	opts := masker.DefaultOptions()
	masker.MaskMap(env, []string{"SECRET"}, opts)
	if env["SECRET"] != "plaintext" {
		t.Error("original map should not be modified")
	}
}
