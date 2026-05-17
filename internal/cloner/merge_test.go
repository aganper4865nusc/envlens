package cloner_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/cloner"
)

func TestCloneInto_NoConflict(t *testing.T) {
	base := makeEnv("A", "1")
	src := makeEnv("B", "2")
	res := cloner.CloneInto(base, src, cloner.DefaultOptions(), cloner.PolicyOverwrite)
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Errorf("unexpected env: %v", res.Env)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestCloneInto_OverwritePolicy(t *testing.T) {
	base := makeEnv("HOST", "old")
	src := makeEnv("HOST", "new")
	res := cloner.CloneInto(base, src, cloner.DefaultOptions(), cloner.PolicyOverwrite)
	if res.Env["HOST"] != "new" {
		t.Errorf("expected overwrite, got %q", res.Env["HOST"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "HOST" {
		t.Errorf("conflict not recorded: %v", res.Conflicts)
	}
}

func TestCloneInto_KeepBasePolicy(t *testing.T) {
	base := makeEnv("HOST", "old")
	src := makeEnv("HOST", "new")
	res := cloner.CloneInto(base, src, cloner.DefaultOptions(), cloner.PolicyKeepBase)
	if res.Env["HOST"] != "old" {
		t.Errorf("expected base value to be kept, got %q", res.Env["HOST"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("conflict should still be recorded: %v", res.Conflicts)
	}
}

func TestCloneInto_WithTransformOnSrc(t *testing.T) {
	base := makeEnv("EXISTING", "yes")
	src := makeEnv("new_key", "  value  ")
	opts := cloner.Options{UppercaseKeys: true, TrimValues: true}
	res := cloner.CloneInto(base, src, opts, cloner.PolicyOverwrite)
	if res.Env["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", res.Env["NEW_KEY"])
	}
	if res.Env["EXISTING"] != "yes" {
		t.Error("base key should be preserved")
	}
}

func TestCloneInto_BaseUnmodified(t *testing.T) {
	base := makeEnv("KEY", "original")
	src := makeEnv("KEY", "override")
	cloner.CloneInto(base, src, cloner.DefaultOptions(), cloner.PolicyOverwrite)
	if base["KEY"] != "original" {
		t.Error("base map was modified")
	}
}
