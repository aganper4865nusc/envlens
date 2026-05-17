package inspector_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/inspector"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestInspect_SourcePreserved(t *testing.T) {
	env := makeEnv("APP_PORT", "8080")
	r := inspector.Inspect("staging.env", env)
	if r.Source != "staging.env" {
		t.Errorf("expected source staging.env, got %s", r.Source)
	}
}

func TestInspect_ProfileTotalKeys(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "")
	r := inspector.Inspect("test", env)
	if r.Profile.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", r.Profile.TotalKeys)
	}
}

func TestInspect_LintIssuesReturned(t *testing.T) {
	// lowercase key should trigger a lint warning
	env := makeEnv("lowercase_key", "value")
	r := inspector.Inspect("test", env)
	if len(r.Lint) == 0 {
		t.Error("expected lint issues for lowercase key, got none")
	}
}

func TestInspect_AuditIssuesReturned(t *testing.T) {
	// empty secret value should be flagged
	env := makeEnv("SECRET_KEY", "")
	r := inspector.Inspect("test", env)
	if len(r.Audit) == 0 {
		t.Error("expected audit issues for empty secret, got none")
	}
}

func TestInspect_ScoreNotNegative(t *testing.T) {
	env := makeEnv("bad key", "", "another bad", "")
	r := inspector.Inspect("test", env)
	if r.Score.Value < 0 {
		t.Errorf("score should not be negative, got %d", r.Score.Value)
	}
}

func TestHasIssues_True(t *testing.T) {
	env := makeEnv("lowercase", "")
	r := inspector.Inspect("test", env)
	if !inspector.HasIssues(r) {
		t.Error("expected HasIssues to return true")
	}
}

func TestHasIssues_False(t *testing.T) {
	env := makeEnv("APP_PORT", "8080", "DB_HOST", "localhost")
	r := inspector.Inspect("test", env)
	if inspector.HasIssues(r) {
		t.Errorf("expected no issues, got lint=%v audit=%v", r.Lint, r.Audit)
	}
}
