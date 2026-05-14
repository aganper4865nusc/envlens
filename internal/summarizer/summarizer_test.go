package summarizer_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/summarizer"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestSummarize_TotalKeys(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "")
	s := summarizer.Summarize(env)
	if s.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", s.TotalKeys)
	}
}

func TestSummarize_EmptyAndNonEmpty(t *testing.T) {
	env := makeEnv("X", "", "Y", "hello", "Z", "world")
	s := summarizer.Summarize(env)
	if s.EmptyValues != 1 {
		t.Errorf("expected 1 empty, got %d", s.EmptyValues)
	}
	if s.NonEmptyValues != 2 {
		t.Errorf("expected 2 non-empty, got %d", s.NonEmptyValues)
	}
}

func TestSummarize_UniqueValues(t *testing.T) {
	env := makeEnv("A", "same", "B", "same", "C", "different")
	s := summarizer.Summarize(env)
	// "same" counted once, "different" counted once => 2 unique
	if s.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", s.UniqueValues)
	}
}

func TestSummarize_LongestKey(t *testing.T) {
	env := makeEnv("SHORT", "v", "MUCH_LONGER_KEY", "v")
	s := summarizer.Summarize(env)
	if s.LongestKey != "MUCH_LONGER_KEY" {
		t.Errorf("expected MUCH_LONGER_KEY as longest, got %s", s.LongestKey)
	}
}

func TestSummarize_Categories(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"JWT_SECRET", "abc",
		"APP_NAME", "envlens",
	)
	s := summarizer.Summarize(env)
	if s.Categories["database"] != 2 {
		t.Errorf("expected 2 database keys, got %d", s.Categories["database"])
	}
	if s.Categories["auth"] < 1 {
		t.Errorf("expected at least 1 auth key, got %d", s.Categories["auth"])
	}
}

func TestSummarize_LinesNotEmpty(t *testing.T) {
	env := makeEnv("FOO", "bar")
	s := summarizer.Summarize(env)
	if len(s.Lines) == 0 {
		t.Error("expected non-empty Lines slice")
	}
}

func TestSummarize_LinesContainTotalKeys(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	s := summarizer.Summarize(env)
	found := false
	for _, l := range s.Lines {
		if strings.Contains(l, "2") && strings.Contains(l, "Total") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Lines did not contain total key count line; got: %v", s.Lines)
	}
}

func TestSummarize_EmptyEnv(t *testing.T) {
	s := summarizer.Summarize(map[string]string{})
	if s.TotalKeys != 0 {
		t.Errorf("expected 0 total keys for empty env, got %d", s.TotalKeys)
	}
	if len(s.Categories) != 0 {
		t.Errorf("expected no categories for empty env")
	}
}
