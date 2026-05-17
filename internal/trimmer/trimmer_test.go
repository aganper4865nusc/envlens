package trimmer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/trimmer"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestTrim_RemoveEmpty(t *testing.T) {
	env := makeEnv("KEY", "value", "EMPTY", "", "ALSO_EMPTY", "")
	res := trimmer.Trim(env, trimmer.Options{RemoveEmpty: true})
	if _, ok := res.Env["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := res.Env["ALSO_EMPTY"]; ok {
		t.Error("expected ALSO_EMPTY to be removed")
	}
	if res.Env["KEY"] != "value" {
		t.Error("expected KEY to be retained")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestTrim_RemoveExplicitKeys(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "3")
	res := trimmer.Trim(env, trimmer.Options{RemoveKeys: []string{"A", "C"}})
	if _, ok := res.Env["A"]; ok {
		t.Error("expected A to be removed")
	}
	if _, ok := res.Env["C"]; ok {
		t.Error("expected C to be removed")
	}
	if res.Env["B"] != "2" {
		t.Error("expected B to be retained")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestTrim_RemoveBlankKeys(t *testing.T) {
	env := map[string]string{"  ": "ghost", "REAL": "yes"}
	res := trimmer.Trim(env, trimmer.Options{RemoveBlankKeys: true})
	if _, ok := res.Env["  "]; ok {
		t.Error("expected blank key to be removed")
	}
	if res.Env["REAL"] != "yes" {
		t.Error("expected REAL to be retained")
	}
}

func TestTrim_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("X", "", "Y", "keep")
	orig := map[string]string{"X": "", "Y": "keep"}
	trimmer.Trim(env, trimmer.Options{RemoveEmpty: true})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("original map was modified at key %q", k)
		}
	}
}

func TestTrim_NoOptions_ReturnsFullCopy(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	res := trimmer.Trim(env, trimmer.Options{})
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
}

func TestTrim_RemoveEmptyAndExplicitKeys_Combined(t *testing.T) {
	env := makeEnv("A", "", "B", "keep", "C", "drop")
	res := trimmer.Trim(env, trimmer.Options{
		RemoveEmpty: true,
		RemoveKeys:  []string{"C"},
	})
	if _, ok := res.Env["A"]; ok {
		t.Error("expected A to be removed (empty value)")
	}
	if _, ok := res.Env["C"]; ok {
		t.Error("expected C to be removed (explicit key)")
	}
	if res.Env["B"] != "keep" {
		t.Error("expected B to be retained")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}
