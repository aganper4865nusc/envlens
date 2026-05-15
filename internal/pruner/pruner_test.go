package pruner

import (
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestPrune_ByExplicitKey(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res, err := Prune(env, Options{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pruned["FOO"]; ok {
		t.Error("expected FOO to be pruned")
	}
	if res.Pruned["BAZ"] != "qux" {
		t.Error("expected BAZ to be retained")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "FOO" {
		t.Errorf("unexpected removed list: %v", res.Removed)
	}
}

func TestPrune_ByPrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	res, err := Prune(env, Options{Prefixes: []string{"DB_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pruned["DB_HOST"]; ok {
		t.Error("expected DB_HOST pruned")
	}
	if _, ok := res.Pruned["DB_PORT"]; ok {
		t.Error("expected DB_PORT pruned")
	}
	if res.Pruned["APP_NAME"] != "envlens" {
		t.Error("expected APP_NAME retained")
	}
}

func TestPrune_ByPattern(t *testing.T) {
	env := makeEnv("SECRET_KEY", "abc", "API_SECRET", "xyz", "HOST", "localhost")
	res, err := Prune(env, Options{Patterns: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pruned["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should be pruned")
	}
	if _, ok := res.Pruned["API_SECRET"]; ok {
		t.Error("API_SECRET should be pruned")
	}
	if res.Pruned["HOST"] != "localhost" {
		t.Error("HOST should be retained")
	}
}

func TestPrune_DropEmpty(t *testing.T) {
	env := makeEnv("EMPTY", "", "NONEMPTY", "value")
	res, err := Prune(env, Options{DropEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Pruned["EMPTY"]; ok {
		t.Error("EMPTY should be pruned")
	}
	if res.Pruned["NONEMPTY"] != "value" {
		t.Error("NONEMPTY should be retained")
	}
}

func TestPrune_InvalidPattern_ReturnsError(t *testing.T) {
	env := makeEnv("FOO", "bar")
	_, err := Prune(env, Options{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestPrune_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	_, err := Prune(env, Options{Keys: []string{"FOO"}, DropEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 2 {
		t.Error("original map was modified")
	}
}

func TestPrune_RemovedListIsSorted(t *testing.T) {
	env := makeEnv("Z_KEY", "1", "A_KEY", "2", "M_KEY", "3")
	res, err := Prune(env, Options{Prefixes: []string{""}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Removed); i++ {
		if res.Removed[i-1] > res.Removed[i] {
			t.Errorf("removed list not sorted: %v", res.Removed)
		}
	}
}
