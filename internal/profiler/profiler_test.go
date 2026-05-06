package profiler_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/profiler"
)

func TestAnalyze_TotalKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	p := profiler.Analyze("test", env)
	if p.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", p.TotalKeys)
	}
	if p.Source != "test" {
		t.Errorf("expected source 'test', got %q", p.Source)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	env := map[string]string{"FOO": "", "BAR": "value"}
	p := profiler.Analyze("s", env)
	if p.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", p.EmptyValues)
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
	p := profiler.Analyze("s", env)
	if len(p.SensitiveKeys) != 2 {
		t.Errorf("expected 2 sensitive keys, got %d: %v", len(p.SensitiveKeys), p.SensitiveKeys)
	}
}

func TestAnalyze_KeyCasing(t *testing.T) {
	env := map[string]string{
		"UPPER_KEY": "a",
		"lower_key": "b",
		"Mixed_Key": "c",
	}
	p := profiler.Analyze("s", env)
	if p.UppercaseKeys != 1 {
		t.Errorf("expected 1 uppercase key, got %d", p.UppercaseKeys)
	}
	if p.LowercaseKeys != 1 {
		t.Errorf("expected 1 lowercase key, got %d", p.LowercaseKeys)
	}
	if p.MixedKeys != 1 {
		t.Errorf("expected 1 mixed key, got %d", p.MixedKeys)
	}
}

func TestAnalyze_UniqueValues(t *testing.T) {
	env := map[string]string{"A": "same", "B": "same", "C": "different"}
	p := profiler.Analyze("s", env)
	if p.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", p.UniqueValues)
	}
}

func TestAnalyze_EmptyMap(t *testing.T) {
	p := profiler.Analyze("empty", map[string]string{})
	if p.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", p.TotalKeys)
	}
	if len(p.SensitiveKeys) != 0 {
		t.Errorf("expected no sensitive keys, got %v", p.SensitiveKeys)
	}
}
