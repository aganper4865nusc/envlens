package deduper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/deduper"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := []deduper.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	r := deduper.Dedupe(entries)
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", r.Duplicates)
	}
	if r.Env["FOO"] != "bar" || r.Env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", r.Env)
	}
}

func TestDedupe_LastValueWins(t *testing.T) {
	entries := []deduper.Entry{
		{Key: "FOO", Value: "first"},
		{Key: "FOO", Value: "second"},
		{Key: "FOO", Value: "third"},
	}
	r := deduper.Dedupe(entries)
	if r.Env["FOO"] != "third" {
		t.Errorf("expected 'third', got %q", r.Env["FOO"])
	}
}

func TestDedupe_DuplicateCountTracked(t *testing.T) {
	entries := []deduper.Entry{
		{Key: "A", Value: "1"},
		{Key: "A", Value: "2"},
		{Key: "B", Value: "x"},
	}
	r := deduper.Dedupe(entries)
	if r.Duplicates["A"] != 2 {
		t.Errorf("expected count 2 for A, got %d", r.Duplicates["A"])
	}
	if _, ok := r.Duplicates["B"]; ok {
		t.Error("B should not appear in duplicates")
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	r := deduper.Dedupe(nil)
	if len(r.Env) != 0 {
		t.Errorf("expected empty env, got %v", r.Env)
	}
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", r.Duplicates)
	}
}

func TestDedupe_MultipleDuplicateKeys(t *testing.T) {
	entries := []deduper.Entry{
		{Key: "X", Value: "a"}, {Key: "X", Value: "b"},
		{Key: "Y", Value: "c"}, {Key: "Y", Value: "d"}, {Key: "Y", Value: "e"},
		{Key: "Z", Value: "z"},
	}
	r := deduper.Dedupe(entries)
	if r.Duplicates["X"] != 2 {
		t.Errorf("X: expected 2, got %d", r.Duplicates["X"])
	}
	if r.Duplicates["Y"] != 3 {
		t.Errorf("Y: expected 3, got %d", r.Duplicates["Y"])
	}
	if r.Env["Y"] != "e" {
		t.Errorf("Y last value: expected 'e', got %q", r.Env["Y"])
	}
}

func TestFromMap_NoFalseDuplicates(t *testing.T) {
	m := map[string]string{"KEY": "val", "OTHER": "x"}
	r := deduper.FromMap(m)
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates from map, got %v", r.Duplicates)
	}
	if r.Env["KEY"] != "val" {
		t.Errorf("unexpected value for KEY: %q", r.Env["KEY"])
	}
}

// TestDedupe_SingleEntry verifies that a single entry produces no duplicates
// and is correctly stored in the result env map.
func TestDedupe_SingleEntry(t *testing.T) {
	entries := []deduper.Entry{
		{Key: "ONLY", Value: "one"},
	}
	r := deduper.Dedupe(entries)
	if len(r.Duplicates) != 0 {
		t.Errorf("expected no duplicates for single entry, got %v", r.Duplicates)
	}
	if r.Env["ONLY"] != "one" {
		t.Errorf("expected 'one', got %q", r.Env["ONLY"])
	}
	if len(r.Env) != 1 {
		t.Errorf("expected env length 1, got %d", len(r.Env))
	}
}
