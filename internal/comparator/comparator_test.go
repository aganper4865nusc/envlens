package comparator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCompare_AllMatch(t *testing.T) {
	left := makeEnv("A", "1", "B", "2")
	right := makeEnv("A", "1", "B", "2")
	res := comparator.Compare(left, right)
	if res.Matched != 2 || res.Total != 2 {
		t.Fatalf("expected 2/2 matched, got %d/%d", res.Matched, res.Total)
	}
	if res.MatchRate != 1.0 {
		t.Fatalf("expected match rate 1.0, got %f", res.MatchRate)
	}
}

func TestCompare_MissingKey(t *testing.T) {
	left := makeEnv("A", "1", "B", "2")
	right := makeEnv("A", "1")
	res := comparator.Compare(left, right)
	found := false
	for _, e := range res.Entries {
		if e.Key == "B" && e.Status == comparator.StatusMissing {
			found = true
		}
	}
	if !found {
		t.Error("expected key B to have StatusMissing")
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	left := makeEnv("A", "1")
	right := makeEnv("A", "1", "C", "3")
	res := comparator.Compare(left, right)
	found := false
	for _, e := range res.Entries {
		if e.Key == "C" && e.Status == comparator.StatusExtra {
			found = true
		}
	}
	if !found {
		t.Error("expected key C to have StatusExtra")
	}
}

func TestCompare_DifferentValue(t *testing.T) {
	left := makeEnv("HOST", "localhost")
	right := makeEnv("HOST", "prod.example.com")
	res := comparator.Compare(left, right)
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Status != comparator.StatusDiffer {
		t.Errorf("expected StatusDiffer, got %s", res.Entries[0].Status)
	}
	if res.Matched != 0 || res.MatchRate != 0.0 {
		t.Errorf("expected 0 matched, got %d / %f", res.Matched, res.MatchRate)
	}
}

func TestCompare_SortedEntries(t *testing.T) {
	left := makeEnv("Z", "z", "A", "a", "M", "m")
	right := makeEnv("Z", "z", "A", "a", "M", "m")
	res := comparator.Compare(left, right)
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %s, got %s", i, k, keys[i])
		}
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := comparator.Compare(map[string]string{}, map[string]string{})
	if res.Total != 0 || res.MatchRate != 0.0 {
		t.Errorf("expected empty result, got total=%d rate=%f", res.Total, res.MatchRate)
	}
}
