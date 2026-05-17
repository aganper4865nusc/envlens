package squasher_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/squasher"
)

func TestSquash_LastWins_NoConflict(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	res := squasher.Squash([]map[string]string{a, b}, squasher.LastWins)
	if res.Env["A"] != "1" || res.Env["B"] != "2" {
		t.Fatalf("unexpected env: %v", res.Env)
	}
	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestSquash_LastWins_ConflictResolved(t *testing.T) {
	a := map[string]string{"KEY": "old"}
	b := map[string]string{"KEY": "new"}
	res := squasher.Squash([]map[string]string{a, b}, squasher.LastWins)
	if res.Env["KEY"] != "new" {
		t.Fatalf("expected 'new', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	c := res.Conflicts[0]
	if c.Kept != "new" || c.Dropped != "old" || c.WonIndex != 1 {
		t.Fatalf("unexpected conflict: %+v", c)
	}
}

func TestSquash_FirstWins_ConflictResolved(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res := squasher.Squash([]map[string]string{a, b}, squasher.FirstWins)
	if res.Env["KEY"] != "first" {
		t.Fatalf("expected 'first', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].WonIndex != 0 {
		t.Fatalf("expected WonIndex 0, got %d", res.Conflicts[0].WonIndex)
	}
}

func TestSquash_ConflictsSorted(t *testing.T) {
	a := map[string]string{"Z": "z1", "A": "a1"}
	b := map[string]string{"Z": "z2", "A": "a2"}
	res := squasher.Squash([]map[string]string{a, b}, squasher.LastWins)
	if len(res.Conflicts) != 2 {
		t.Fatalf("expected 2 conflicts, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "A" || res.Conflicts[1].Key != "Z" {
		t.Fatalf("conflicts not sorted: %v", res.Conflicts)
	}
}

func TestSquash_EmptySources(t *testing.T) {
	res := squasher.Squash([]map[string]string{}, squasher.LastWins)
	if len(res.Env) != 0 {
		t.Fatalf("expected empty env")
	}
}

func TestSquash_ThreeSources_LastWins(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	c := map[string]string{"X": "3"}
	res := squasher.Squash([]map[string]string{a, b, c}, squasher.LastWins)
	if res.Env["X"] != "3" {
		t.Fatalf("expected '3', got %q", res.Env["X"])
	}
	// Two conflicts: a vs b, then (b) vs c
	if len(res.Conflicts) != 2 {
		t.Fatalf("expected 2 conflicts, got %d", len(res.Conflicts))
	}
}
