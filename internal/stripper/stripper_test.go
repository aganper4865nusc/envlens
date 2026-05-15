package stripper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/stripper"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestStrip_ByPrefix(t *testing.T) {
	env := makeEnv("INTERNAL_FOO", "1", "INTERNAL_BAR", "2", "PUBLIC_KEY", "3")
	res, err := stripper.Strip(env, stripper.Options{Prefixes: []string{"INTERNAL_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["PUBLIC_KEY"]; !ok {
		t.Error("expected PUBLIC_KEY to remain")
	}
	if len(res.Env) != 1 {
		t.Errorf("expected 1 key, got %d", len(res.Env))
	}
	if len(res.Stripped) != 2 {
		t.Errorf("expected 2 stripped keys, got %d", len(res.Stripped))
	}
}

func TestStrip_ByExplicitKeys(t *testing.T) {
	env := makeEnv("FOO", "a", "BAR", "b", "BAZ", "c")
	res, err := stripper.Strip(env, stripper.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["BAR"]; !ok {
		t.Error("expected BAR to remain")
	}
	if len(res.Stripped) != 2 {
		t.Errorf("expected 2 stripped, got %d", len(res.Stripped))
	}
}

func TestStrip_ByPattern(t *testing.T) {
	env := makeEnv("SECRET_KEY", "x", "SECRET_TOKEN", "y", "APP_NAME", "z")
	res, err := stripper.Strip(env, stripper.Options{Pattern: "^SECRET_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 1 {
		t.Errorf("expected 1 key remaining, got %d", len(res.Env))
	}
}

func TestStrip_InvalidPattern_ReturnsError(t *testing.T) {
	env := makeEnv("FOO", "bar")
	_, err := stripper.Strip(env, stripper.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestStrip_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	_, err := stripper.Strip(env, stripper.Options{Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["A"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestStrip_StrippedKeysSorted(t *testing.T) {
	env := makeEnv("Z_KEY", "1", "A_KEY", "2", "M_KEY", "3")
	res, err := stripper.Strip(env, stripper.Options{Prefixes: []string{""}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Stripped); i++ {
		if res.Stripped[i] < res.Stripped[i-1] {
			t.Errorf("stripped keys not sorted: %v", res.Stripped)
			break
		}
	}
}

func TestStrip_NoOptions_NothingRemoved(t *testing.T) {
	env := makeEnv("FOO", "1", "BAR", "2")
	res, err := stripper.Strip(env, stripper.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if len(res.Stripped) != 0 {
		t.Errorf("expected no stripped keys, got %v", res.Stripped)
	}
}
