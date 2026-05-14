package scoper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/scoper"
)

func TestMerge_NoConflict(t *testing.T) {
	r1 := scoper.Result{Scope: "a", Env: map[string]string{"HOST": "h1"}}
	r2 := scoper.Result{Scope: "b", Env: map[string]string{"PORT": "5432"}}
	m := scoper.Merge([]scoper.Result{r1, r2}, scoper.MergeOptions{})
	if len(m.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(m.Env))
	}
	if len(m.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", m.Conflicts)
	}
}

func TestMerge_ConflictRecorded(t *testing.T) {
	r1 := scoper.Result{Scope: "a", Env: map[string]string{"HOST": "h1"}}
	r2 := scoper.Result{Scope: "b", Env: map[string]string{"HOST": "h2"}}
	m := scoper.Merge([]scoper.Result{r1, r2}, scoper.MergeOptions{Overwrite: false})
	if len(m.Conflicts) != 1 || m.Conflicts[0] != "HOST" {
		t.Fatalf("expected conflict on HOST, got %v", m.Conflicts)
	}
	if m.Env["HOST"] != "h1" {
		t.Fatalf("expected first value to win, got %q", m.Env["HOST"])
	}
}

func TestMerge_OverwriteEnabled(t *testing.T) {
	r1 := scoper.Result{Scope: "a", Env: map[string]string{"HOST": "h1"}}
	r2 := scoper.Result{Scope: "b", Env: map[string]string{"HOST": "h2"}}
	m := scoper.Merge([]scoper.Result{r1, r2}, scoper.MergeOptions{Overwrite: true})
	if len(m.Conflicts) != 0 {
		t.Fatalf("expected no conflicts when overwrite=true, got %v", m.Conflicts)
	}
	if m.Env["HOST"] != "h2" {
		t.Fatalf("expected last value to win, got %q", m.Env["HOST"])
	}
}

func TestMerge_PrefixApplied(t *testing.T) {
	r := scoper.Result{Scope: "x", Env: map[string]string{"KEY": "val"}}
	m := scoper.Merge([]scoper.Result{r}, scoper.MergeOptions{Prefix: "APP_"})
	if _, ok := m.Env["APP_KEY"]; !ok {
		t.Fatal("expected APP_KEY in merged output")
	}
	if _, ok := m.Env["KEY"]; ok {
		t.Fatal("unprefixed KEY should not exist")
	}
}

func TestMerge_EmptyResults(t *testing.T) {
	m := scoper.Merge([]scoper.Result{}, scoper.MergeOptions{})
	if len(m.Env) != 0 {
		t.Fatalf("expected empty env, got %d keys", len(m.Env))
	}
}
