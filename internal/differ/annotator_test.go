package differ

import (
	"testing"
)

func makeEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "A", Status: "added", OldValue: "", NewValue: "1"},
		{Key: "B", Status: "removed", OldValue: "2", NewValue: ""},
		{Key: "C", Status: "changed", OldValue: "old", NewValue: "new"},
		{Key: "D", Status: "unchanged", OldValue: "same", NewValue: "same"},
	}
}

func TestAnnotate_StatusTagApplied(t *testing.T) {
	results := Annotate(makeEntries(), DefaultAnnotateOptions())

	expected := map[string]string{
		"A": "added",
		"B": "removed",
		"C": "changed",
		"D": "unchanged",
	}

	for _, a := range results {
		want, ok := expected[a.Key]
		if !ok {
			t.Fatalf("unexpected key %q", a.Key)
		}
		if a.Status != want {
			t.Errorf("key %q: got status %q, want %q", a.Key, a.Status, want)
		}
		found := false
		for _, tag := range a.Tags {
			if tag == want {
				found = true
			}
		}
		if !found {
			t.Errorf("key %q: tag %q not found in %v", a.Key, want, a.Tags)
		}
	}
}

func TestAnnotate_ExtraTagsApplied(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.ExtraTags = []string{"env:prod", "reviewed"}

	results := Annotate(makeEntries(), opts)
	for _, a := range results {
		hasEnv, hasReviewed := false, false
		for _, tag := range a.Tags {
			if tag == "env:prod" {
				hasEnv = true
			}
			if tag == "reviewed" {
				hasReviewed = true
			}
		}
		if !hasEnv || !hasReviewed {
			t.Errorf("key %q missing extra tags, got %v", a.Key, a.Tags)
		}
	}
}

func TestAnnotate_ValuesPreserved(t *testing.T) {
	results := Annotate(makeEntries(), DefaultAnnotateOptions())
	for _, a := range results {
		if a.Key == "C" {
			if a.OldValue != "old" || a.NewValue != "new" {
				t.Errorf("values not preserved: old=%q new=%q", a.OldValue, a.NewValue)
			}
		}
	}
}

func TestFilterAnnotations_ByStatus(t *testing.T) {
	results := Annotate(makeEntries(), DefaultAnnotateOptions())
	filtered := FilterAnnotations(results, "added", "removed")

	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered, got %d", len(filtered))
	}
	for _, a := range filtered {
		if a.Status != "added" && a.Status != "removed" {
			t.Errorf("unexpected status %q in filtered results", a.Status)
		}
	}
}

func TestFilterAnnotations_EmptyStatuses(t *testing.T) {
	results := Annotate(makeEntries(), DefaultAnnotateOptions())
	filtered := FilterAnnotations(results)
	if len(filtered) != 0 {
		t.Errorf("expected 0 results with no statuses, got %d", len(filtered))
	}
}

func TestAnnotate_EmptyInput(t *testing.T) {
	results := Annotate([]DiffEntry{}, DefaultAnnotateOptions())
	if len(results) != 0 {
		t.Errorf("expected empty output, got %d items", len(results))
	}
}
