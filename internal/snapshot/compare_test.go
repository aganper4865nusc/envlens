package snapshot_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/snapshot"
)

func makeSnap(source string, env map[string]string) *snapshot.Snapshot {
	return snapshot.Take(source, env)
}

func TestCompare_Added(t *testing.T) {
	old := makeSnap("old", map[string]string{"A": "1"})
	new := makeSnap("new", map[string]string{"A": "1", "B": "2"})
	d := snapshot.Compare(old, new)
	if _, ok := d.Added["B"]; !ok {
		t.Error("expected B to be in Added")
	}
	if d.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Removed(t *testing.T) {
	old := makeSnap("old", map[string]string{"A": "1", "B": "2"})
	new := makeSnap("new", map[string]string{"A": "1"})
	d := snapshot.Compare(old, new)
	if _, ok := d.Removed["B"]; !ok {
		t.Error("expected B to be in Removed")
	}
}

func TestCompare_Changed(t *testing.T) {
	old := makeSnap("old", map[string]string{"A": "1"})
	new := makeSnap("new", map[string]string{"A": "2"})
	d := snapshot.Compare(old, new)
	pair, ok := d.Changed["A"]
	if !ok {
		t.Fatal("expected A to be in Changed")
	}
	if pair[0] != "1" || pair[1] != "2" {
		t.Errorf("expected [1 2], got %v", pair)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := makeSnap("old", map[string]string{"A": "1", "B": "2"})
	new := makeSnap("new", map[string]string{"A": "1", "B": "2"})
	d := snapshot.Compare(old, new)
	if d.HasChanges() {
		t.Error("expected no changes")
	}
	if len(d.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(d.Unchanged))
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	old := makeSnap("old", map[string]string{})
	new := makeSnap("new", map[string]string{})
	d := snapshot.Compare(old, new)
	if d.HasChanges() {
		t.Error("expected no changes for empty snapshots")
	}
}
