package cloner_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/cloner"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestClone_BasicCopy(t *testing.T) {
	src := makeEnv("FOO", "bar", "BAZ", "qux")
	res := cloner.Clone(src, cloner.DefaultOptions())
	if len(res.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(res.Env))
	}
	if res.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Env["FOO"])
	}
}

func TestClone_OriginalUnmodified(t *testing.T) {
	src := makeEnv("KEY", "value")
	res := cloner.Clone(src, cloner.Options{TrimValues: true})
	res.Env["KEY"] = "mutated"
	if src["KEY"] != "value" {
		t.Error("original map was modified")
	}
}

func TestClone_UppercaseKeys(t *testing.T) {
	src := makeEnv("db_host", "localhost", "db_port", "5432")
	res := cloner.Clone(src, cloner.Options{UppercaseKeys: true})
	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in cloned map")
	}
	if _, ok := res.Env["db_host"]; ok {
		t.Error("lowercase key should not be present")
	}
	if res.TransformedKeys["db_host"] != "DB_HOST" {
		t.Errorf("TransformedKeys mismatch: %v", res.TransformedKeys)
	}
}

func TestClone_TrimValues(t *testing.T) {
	src := makeEnv("HOST", "  localhost  ")
	res := cloner.Clone(src, cloner.Options{TrimValues: true})
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", res.Env["HOST"])
	}
}

func TestClone_ExcludeKeys(t *testing.T) {
	src := makeEnv("SECRET", "s3cr3t", "HOST", "localhost")
	res := cloner.Clone(src, cloner.Options{ExcludeKeys: []string{"SECRET"}})
	if _, ok := res.Env["SECRET"]; ok {
		t.Error("excluded key should not be in cloned map")
	}
	if len(res.ExcludedKeys) != 1 || res.ExcludedKeys[0] != "SECRET" {
		t.Errorf("ExcludedKeys not tracked correctly: %v", res.ExcludedKeys)
	}
}

func TestClone_EmptySource(t *testing.T) {
	res := cloner.Clone(map[string]string{}, cloner.DefaultOptions())
	if len(res.Env) != 0 {
		t.Errorf("expected empty map, got %d keys", len(res.Env))
	}
}

func TestClone_NoTransformKeys_WhenAlreadyUpper(t *testing.T) {
	src := makeEnv("ALREADY", "up")
	res := cloner.Clone(src, cloner.Options{UppercaseKeys: true})
	if len(res.TransformedKeys) != 0 {
		t.Errorf("expected no transformed keys, got %v", res.TransformedKeys)
	}
}
