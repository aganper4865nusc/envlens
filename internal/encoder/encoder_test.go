package encoder_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/encoder"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestEncode_Base64(t *testing.T) {
	env := makeEnv("SECRET", "hello world")
	res, err := encoder.Encode(env, encoder.FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := res.Encoded["SECRET"], "aGVsbG8gd29ybGQ="; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", res.Skipped)
	}
}

func TestEncode_Hex(t *testing.T) {
	env := makeEnv("TOKEN", "abc")
	res, err := encoder.Encode(env, encoder.FormatHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := res.Encoded["TOKEN"], "616263"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncode_URL(t *testing.T) {
	env := makeEnv("REDIRECT", "https://example.com/path?foo=bar&baz=qux")
	res, err := encoder.Encode(env, encoder.FormatURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Encoded["REDIRECT"] == env["REDIRECT"] {
		t.Error("expected value to be URL-encoded")
	}
}

func TestDecode_Base64_RoundTrip(t *testing.T) {
	original := makeEnv("KEY", "supersecret")
	encoded, _ := encoder.Encode(original, encoder.FormatBase64)
	decoded, err := encoder.Decode(encoded.Encoded, encoder.FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := decoded.Encoded["KEY"], "supersecret"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDecode_Hex_RoundTrip(t *testing.T) {
	original := makeEnv("API_KEY", "mytoken123")
	encoded, _ := encoder.Encode(original, encoder.FormatHex)
	decoded, err := encoder.Decode(encoded.Encoded, encoder.FormatHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := decoded.Encoded["API_KEY"], "mytoken123"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDecode_InvalidBase64_Skipped(t *testing.T) {
	env := makeEnv("BAD", "not-valid-base64!!!", "GOOD", "aGVsbG8=")
	res, err := encoder.Decode(env, encoder.FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "BAD" {
		t.Errorf("expected BAD to be skipped, got %v", res.Skipped)
	}
	if got, want := res.Encoded["GOOD"], "hello"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEncode_UnsupportedFormat(t *testing.T) {
	env := makeEnv("X", "y")
	res, err := encoder.Encode(env, encoder.Format("rot13"))
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key for unsupported format, got %d", len(res.Skipped))
	}
}
