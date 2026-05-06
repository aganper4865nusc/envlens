package renamer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/renamer"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestRename_BasicApplied(t *testing.T) {
	env := makeEnv("OLD_KEY", "value")
	res := renamer.Rename(env, []renamer.Rule{{From: "OLD_KEY", To: "NEW_KEY"}})

	if _, ok := res.Env["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to exist")
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if len(res.Applied) != 1 || res.Applied[0].From != "OLD_KEY" {
		t.Errorf("unexpected Applied: %v", res.Applied)
	}
}

func TestRename_SkippedWhenFromMissing(t *testing.T) {
	env := makeEnv("A", "1")
	res := renamer.Rename(env, []renamer.Rule{{From: "MISSING", To: "B"}})

	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Applied) != 0 {
		t.Errorf("expected 0 applied, got %d", len(res.Applied))
	}
}

func TestRename_ConflictWhenDestExists(t *testing.T) {
	env := makeEnv("SRC", "hello", "DST", "world")
	res := renamer.Rename(env, []renamer.Rule{{From: "SRC", To: "DST"}})

	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	// Original map must be untouched
	if res.Env["DST"] != "world" {
		t.Errorf("DST should retain original value, got %q", res.Env["DST"])
	}
}

func TestRename_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("K", "v")
	renamer.Rename(env, []renamer.Rule{{From: "K", To: "K2"}})

	if _, ok := env["K"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestRename_MultipleRules(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "3")
	rules := []renamer.Rule{
		{From: "A", To: "ALPHA"},
		{From: "B", To: "BETA"},
		{From: "NOPE", To: "X"},
	}
	res := renamer.Rename(env, rules)

	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if res.Env["ALPHA"] != "1" || res.Env["BETA"] != "2" || res.Env["C"] != "3" {
		t.Errorf("unexpected env state: %v", res.Env)
	}
}

func TestRename_TrimSpaceInRuleKeys(t *testing.T) {
	env := makeEnv("MYKEY", "val")
	res := renamer.Rename(env, []renamer.Rule{{From: "  MYKEY  ", To: "  NEWKEY  "}})

	if _, ok := res.Env["NEWKEY"]; !ok {
		t.Error("expected NEWKEY after trimming rule keys")
	}
}
