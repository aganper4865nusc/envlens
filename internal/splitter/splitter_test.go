package splitter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/splitter"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestSplit_ByPrefix(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens", "LOG_LEVEL", "info")
	rules := []splitter.Rule{
		{Name: "database", Prefix: "DB_"},
		{Name: "app", Prefix: "APP_"},
	}
	res, err := splitter.Split(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Buckets["database"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST in database bucket")
	}
	if res.Buckets["app"]["APP_NAME"] != "envlens" {
		t.Errorf("expected APP_NAME in app bucket")
	}
	if _, ok := res.Unmatched["LOG_LEVEL"]; !ok {
		t.Errorf("expected LOG_LEVEL in unmatched")
	}
}

func TestSplit_ByPattern(t *testing.T) {
	env := makeEnv("AWS_KEY", "k1", "AWS_SECRET", "s1", "GCP_TOKEN", "t1")
	rules := []splitter.Rule{
		{Name: "cloud", Pattern: `^(AWS|GCP)_`},
	}
	res, err := splitter.Split(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["cloud"]) != 3 {
		t.Errorf("expected 3 keys in cloud bucket, got %d", len(res.Buckets["cloud"]))
	}
	if len(res.Unmatched) != 0 {
		t.Errorf("expected no unmatched keys")
	}
}

func TestSplit_FirstRuleWins(t *testing.T) {
	env := makeEnv("DB_URL", "postgres://")
	rules := []splitter.Rule{
		{Name: "first", Prefix: "DB_"},
		{Name: "second", Prefix: "DB_"},
	}
	res, err := splitter.Split(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Buckets["first"]["DB_URL"]; !ok {
		t.Errorf("expected DB_URL in first bucket")
	}
	if _, ok := res.Buckets["second"]["DB_URL"]; ok {
		t.Errorf("DB_URL should not appear in second bucket")
	}
}

func TestSplit_InvalidPattern_ReturnsError(t *testing.T) {
	env := makeEnv("KEY", "val")
	rules := []splitter.Rule{
		{Name: "bad", Pattern: `[invalid`},
	}
	_, err := splitter.Split(env, rules)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestSplit_EmptyEnv(t *testing.T) {
	res, err := splitter.Split(map[string]string{}, []splitter.Rule{
		{Name: "db", Prefix: "DB_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["db"]) != 0 {
		t.Errorf("expected empty db bucket")
	}
	if len(res.Unmatched) != 0 {
		t.Errorf("expected empty unmatched")
	}
}
