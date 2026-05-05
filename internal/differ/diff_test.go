package differ_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW": "val"}

	entries := differ.Diff(base, target)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	found := findEntry(entries, "NEW")
	if found == nil {
		t.Fatal("expected entry for NEW")
	}
	if found.Status != differ.StatusAdded {
		t.Errorf("expected added, got %s", found.Status)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	target := map[string]string{"FOO": "bar"}

	entries := differ.Diff(base, target)
	found := findEntry(entries, "OLD")
	if found == nil {
		t.Fatal("expected entry for OLD")
	}
	if found.Status != differ.StatusRemoved {
		t.Errorf("expected removed, got %s", found.Status)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	base := map[string]string{"DB_URL": "localhost"}
	target := map[string]string{"DB_URL": "prod.db.internal"}

	entries := differ.Diff(base, target)
	found := findEntry(entries, "DB_URL")
	if found == nil {
		t.Fatal("expected entry for DB_URL")
	}
	if found.Status != differ.StatusChanged {
		t.Errorf("expected changed, got %s", found.Status)
	}
	if found.BaseVal != "localhost" || found.TargetVal != "prod.db.internal" {
		t.Errorf("unexpected values: %+v", found)
	}
}

func TestDiff_SameKey(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	target := map[string]string{"PORT": "8080"}

	entries := differ.Diff(base, target)
	found := findEntry(entries, "PORT")
	if found == nil || found.Status != differ.StatusSame {
		t.Errorf("expected same status for PORT")
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	target := map[string]string{"Z": "1", "A": "2", "M": "3"}

	entries := differ.Diff(base, target)
	keys := []string{entries[0].Key, entries[1].Key, entries[2].Key}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys A,M,Z got %v", keys)
	}
}

func findEntry(entries []differ.DiffEntry, key string) *differ.DiffEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
