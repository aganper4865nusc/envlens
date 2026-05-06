package filter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/filter"
)

func makeEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.internal",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result, err := filter.Filter(makeEnv(), filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched keys, got %d", len(result.Matched))
	}
	if _, ok := result.Matched["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in matched")
	}
	if _, ok := result.Matched["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in matched")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	result, err := filter.Filter(makeEnv(), filter.Options{Pattern: ".*PASSWORD.*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 1 {
		t.Errorf("expected 1 matched key, got %d", len(result.Matched))
	}
	if _, ok := result.Matched["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in matched")
	}
}

func TestFilter_Inverted(t *testing.T) {
	result, err := filter.Filter(makeEnv(), filter.Options{Prefix: "DB_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != 3 {
		t.Errorf("expected 3 matched keys, got %d", len(result.Matched))
	}
	if _, ok := result.Excluded["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in excluded")
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.Filter(makeEnv(), filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_NoOptions_AllMatched(t *testing.T) {
	env := makeEnv()
	result, err := filter.Filter(env, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Matched) != len(env) {
		t.Errorf("expected all %d keys matched, got %d", len(env), len(result.Matched))
	}
	if len(result.Excluded) != 0 {
		t.Errorf("expected 0 excluded keys, got %d", len(result.Excluded))
	}
}
