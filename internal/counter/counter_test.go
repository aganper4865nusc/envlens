package counter_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/counter"
)

func TestCount_TotalKeys(t *testing.T) {
	srcs := []map[string]string{
		{"A": "1", "B": "2"},
		{"A": "1", "C": "3"},
	}
	r := counter.Count(srcs)
	if r.TotalKeys != 3 {
		t.Fatalf("expected 3 total keys, got %d", r.TotalKeys)
	}
}

func TestCount_OccurrenceTracked(t *testing.T) {
	srcs := []map[string]string{
		{"X": "hello"},
		{"X": "world"},
	}
	r := counter.Count(srcs)
	if len(r.Stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(r.Stats))
	}
	if r.Stats[0].Occurrence != 2 {
		t.Errorf("expected occurrence 2, got %d", r.Stats[0].Occurrence)
	}
}

func TestCount_EmptyValuesTracked(t *testing.T) {
	srcs := []map[string]string{
		{"KEY": ""},
		{"KEY": ""},
		{"KEY": "val"},
	}
	r := counter.Count(srcs)
	if r.TotalEmpty != 2 {
		t.Errorf("expected TotalEmpty=2, got %d", r.TotalEmpty)
	}
	if r.Stats[0].EmptyCount != 2 {
		t.Errorf("expected EmptyCount=2 for KEY, got %d", r.Stats[0].EmptyCount)
	}
}

func TestCount_UniqueValues(t *testing.T) {
	srcs := []map[string]string{
		{"DB": "postgres"},
		{"DB": "mysql"},
		{"DB": "postgres"},
	}
	r := counter.Count(srcs)
	if r.Stats[0].UniqueVals != 2 {
		t.Errorf("expected 2 unique values, got %d", r.Stats[0].UniqueVals)
	}
}

func TestCount_StatsSortedAlphabetically(t *testing.T) {
	srcs := []map[string]string{
		{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"},
	}
	r := counter.Count(srcs)
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, s := range r.Stats {
		if s.Key != expected[i] {
			t.Errorf("pos %d: expected %s, got %s", i, expected[i], s.Key)
		}
	}
}

func TestCount_EmptySources(t *testing.T) {
	r := counter.Count(nil)
	if r.TotalKeys != 0 || r.TotalEmpty != 0 {
		t.Errorf("expected zero result for nil input, got %+v", r)
	}
}
