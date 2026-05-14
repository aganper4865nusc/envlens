package highlighter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/highlighter"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestHighlight_AddedKey(t *testing.T) {
	base := makeEnv("APP_PORT", "8080")
	curr := makeEnv("APP_PORT", "8080", "APP_DEBUG", "true")
	res := highlighter.Highlight(base, curr)
	if res.Added != 1 {
		t.Fatalf("expected 1 added, got %d", res.Added)
	}
	if res.Entries[0].Key != "APP_DEBUG" || res.Entries[0].Status != highlighter.Added {
		t.Errorf("unexpected entry: %+v", res.Entries[0])
	}
}

func TestHighlight_RemovedKey(t *testing.T) {
	base := makeEnv("APP_PORT", "8080", "APP_SECRET", "s3cr3t")
	curr := makeEnv("APP_PORT", "8080")
	res := highlighter.Highlight(base, curr)
	if res.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", res.Removed)
	}
	for _, e := range res.Entries {
		if e.Key == "APP_SECRET" && e.Status != highlighter.Removed {
			t.Errorf("expected APP_SECRET to be removed, got %s", e.Status)
		}
	}
}

func TestHighlight_ChangedKey(t *testing.T) {
	base := makeEnv("DB_HOST", "localhost")
	curr := makeEnv("DB_HOST", "prod.db.internal")
	res := highlighter.Highlight(base, curr)
	if res.Changed != 1 {
		t.Fatalf("expected 1 changed, got %d", res.Changed)
	}
	e := res.Entries[0]
	if e.OldValue != "localhost" || e.Value != "prod.db.internal" {
		t.Errorf("unexpected values: old=%s new=%s", e.OldValue, e.Value)
	}
}

func TestHighlight_UnchangedKey(t *testing.T) {
	base := makeEnv("LOG_LEVEL", "info")
	curr := makeEnv("LOG_LEVEL", "info")
	res := highlighter.Highlight(base, curr)
	if res.Unchanged != 1 || res.Added != 0 || res.Changed != 0 || res.Removed != 0 {
		t.Errorf("unexpected counts: %+v", res)
	}
}

func TestHighlight_SortedOutput(t *testing.T) {
	base := makeEnv("Z_KEY", "1", "A_KEY", "2")
	curr := makeEnv("Z_KEY", "1", "A_KEY", "2", "M_KEY", "3")
	res := highlighter.Highlight(base, curr)
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: want %s got %s", i, k, keys[i])
		}
	}
}

func TestHighlight_EmptyBaseline(t *testing.T) {
	curr := makeEnv("FOO", "bar", "BAZ", "qux")
	res := highlighter.Highlight(map[string]string{}, curr)
	if res.Added != 2 {
		t.Errorf("expected 2 added, got %d", res.Added)
	}
}
