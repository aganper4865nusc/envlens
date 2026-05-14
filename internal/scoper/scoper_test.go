package scoper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/scoper"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestScope_MatchesByPrefix(t *testing.T) {
	env := makeEnv("PROD_DB", "db1", "STAGING_DB", "db2", "COMMON", "val")
	r := scoper.Scope(env, scoper.Rule{Scope: "prod", Prefix: "PROD_"})
	if _, ok := r.Env["PROD_DB"]; !ok {
		t.Fatal("expected PROD_DB in result")
	}
	if _, ok := r.Env["STAGING_DB"]; ok {
		t.Fatal("STAGING_DB should not be in prod scope")
	}
	if r.Scope != "prod" {
		t.Fatalf("expected scope 'prod', got %q", r.Scope)
	}
}

func TestScope_StripPrefix(t *testing.T) {
	env := makeEnv("PROD_HOST", "localhost", "PROD_PORT", "5432")
	r := scoper.Scope(env, scoper.Rule{Scope: "prod", Prefix: "PROD_", StripPrefix: true})
	if v, ok := r.Env["HOST"]; !ok || v != "localhost" {
		t.Fatalf("expected HOST=localhost, got %v/%v", ok, v)
	}
	if _, ok := r.Env["PROD_HOST"]; ok {
		t.Fatal("prefixed key should have been stripped")
	}
	if !r.Stripped {
		t.Fatal("expected Stripped=true")
	}
}

func TestScope_EmptyPrefix_MatchesAll(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	r := scoper.Scope(env, scoper.Rule{Scope: "all", Prefix: ""})
	if len(r.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Env))
	}
}

func TestScope_MatchedKeysSorted(t *testing.T) {
	env := makeEnv("PROD_Z", "z", "PROD_A", "a", "PROD_M", "m")
	r := scoper.Scope(env, scoper.Rule{Scope: "prod", Prefix: "PROD_"})
	for i := 1; i < len(r.MatchedKeys); i++ {
		if r.MatchedKeys[i-1] > r.MatchedKeys[i] {
			t.Fatalf("matched keys not sorted: %v", r.MatchedKeys)
		}
	}
}

func TestScope_NoMatch_EmptyResult(t *testing.T) {
	env := makeEnv("STAGING_X", "1")
	r := scoper.Scope(env, scoper.Rule{Scope: "prod", Prefix: "PROD_"})
	if len(r.Env) != 0 {
		t.Fatalf("expected empty env, got %d keys", len(r.Env))
	}
	if len(r.MatchedKeys) != 0 {
		t.Fatal("expected no matched keys")
	}
}

func TestScopeAll_MultipleRules(t *testing.T) {
	env := makeEnv("PROD_A", "1", "STAGING_A", "2", "COMMON", "3")
	rules := []scoper.Rule{
		{Scope: "prod", Prefix: "PROD_", StripPrefix: true},
		{Scope: "staging", Prefix: "STAGING_", StripPrefix: true},
	}
	results := scoper.ScopeAll(env, rules)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if _, ok := results[0].Env["A"]; !ok {
		t.Fatal("expected 'A' in prod scope after strip")
	}
	if _, ok := results[1].Env["A"]; !ok {
		t.Fatal("expected 'A' in staging scope after strip")
	}
}
