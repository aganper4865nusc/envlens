package grouper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/grouper"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestGroupBy_SinglePrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	res := grouper.GroupBy(env, []string{"DB_"})

	if len(res.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(res.Groups))
	}
	g := res.Groups[0]
	if g.Name != "DB_" {
		t.Errorf("expected group name DB_, got %s", g.Name)
	}
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 keys in DB_ group, got %d", len(g.Keys))
	}
	if _, ok := res.Ungrouped["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in ungrouped")
	}
}

func TestGroupBy_MultiplePrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "AWS_KEY", "abc", "APP_PORT", "8080", "UNRELATED", "x")
	res := grouper.GroupBy(env, []string{"DB_", "AWS_", "APP_"})

	if len(res.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(res.Groups))
	}
	if _, ok := res.Ungrouped["UNRELATED"]; !ok {
		t.Error("expected UNRELATED in ungrouped")
	}
	if len(res.Ungrouped) != 1 {
		t.Errorf("expected 1 ungrouped key, got %d", len(res.Ungrouped))
	}
}

func TestGroupBy_KeysSorted(t *testing.T) {
	env := makeEnv("DB_PORT", "5432", "DB_HOST", "localhost", "DB_NAME", "mydb")
	res := grouper.GroupBy(env, []string{"DB_"})

	keys := res.Groups[0].Keys
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_NAME" || keys[2] != "DB_PORT" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestGroupBy_EmptyEnv(t *testing.T) {
	res := grouper.GroupBy(map[string]string{}, []string{"DB_"})
	if len(res.Groups[0].Keys) != 0 {
		t.Error("expected no keys in group for empty env")
	}
	if len(res.Ungrouped) != 0 {
		t.Error("expected no ungrouped keys for empty env")
	}
}

func TestGroupBy_NoPrefixes(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res := grouper.GroupBy(env, []string{})

	if len(res.Groups) != 0 {
		t.Error("expected no groups")
	}
	if len(res.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(res.Ungrouped))
	}
}

func TestGroupBy_CaseInsensitivePrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	res := grouper.GroupBy(env, []string{"db_"})

	if len(res.Groups[0].Keys) != 1 {
		t.Errorf("expected 1 key matched case-insensitively, got %d", len(res.Groups[0].Keys))
	}
}
