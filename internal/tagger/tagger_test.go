package tagger_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/tagger"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestTag_PrefixMatch(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	rules := []tagger.Rule{{Tag: "database", Prefixes: []string{"DB_"}}}
	results := tagger.Tag(env, rules)
	tagged := map[string][]string{}
	for _, r := range results {
		tagged[r.Key] = r.Tags
	}
	if !containsTag(tagged["DB_HOST"], "database") {
		t.Errorf("expected DB_HOST to have tag 'database'")
	}
	if !containsTag(tagged["DB_PORT"], "database") {
		t.Errorf("expected DB_PORT to have tag 'database'")
	}
	if containsTag(tagged["APP_NAME"], "database") {
		t.Errorf("APP_NAME should not have tag 'database'")
	}
}

func TestTag_ExplicitKeyMatch(t *testing.T) {
	env := makeEnv("SECRET_KEY", "abc", "PUBLIC_URL", "https://x")
	rules := []tagger.Rule{{Tag: "secret", Keys: []string{"SECRET_KEY"}}}
	results := tagger.Tag(env, rules)
	for _, r := range results {
		if r.Key == "SECRET_KEY" && !containsTag(r.Tags, "secret") {
			t.Errorf("expected SECRET_KEY to be tagged 'secret'")
		}
		if r.Key == "PUBLIC_URL" && containsTag(r.Tags, "secret") {
			t.Errorf("PUBLIC_URL should not be tagged 'secret'")
		}
	}
}

func TestTag_PatternMatch(t *testing.T) {
	env := makeEnv("AWS_ACCESS_KEY_ID", "id", "AWS_SECRET", "s", "LOG_LEVEL", "info")
	rules := []tagger.Rule{{Tag: "aws", Pattern: `^AWS_`}}
	results := tagger.Tag(env, rules)
	for _, r := range results {
		hasTag := containsTag(r.Tags, "aws")
		if (r.Key == "AWS_ACCESS_KEY_ID" || r.Key == "AWS_SECRET") && !hasTag {
			t.Errorf("%s should have tag 'aws'", r.Key)
		}
		if r.Key == "LOG_LEVEL" && hasTag {
			t.Errorf("LOG_LEVEL should not have tag 'aws'")
		}
	}
}

func TestTag_MultipleTagsPerKey(t *testing.T) {
	env := makeEnv("DB_PASSWORD", "secret")
	rules := []tagger.Rule{
		{Tag: "database", Prefixes: []string{"DB_"}},
		{Tag: "sensitive", Pattern: `PASSWORD`},
	}
	results := tagger.Tag(env, rules)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !containsTag(results[0].Tags, "database") || !containsTag(results[0].Tags, "sensitive") {
		t.Errorf("expected both tags on DB_PASSWORD, got %v", results[0].Tags)
	}
}

func TestTag_NoDuplicateTags(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	rules := []tagger.Rule{
		{Tag: "infra", Prefixes: []string{"DB_"}},
		{Tag: "infra", Keys: []string{"DB_HOST"}},
	}
	results := tagger.Tag(env, rules)
	count := 0
	for _, tag := range results[0].Tags {
		if tag == "infra" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected tag 'infra' exactly once, got %d", count)
	}
}

func TestTag_EmptyEnv(t *testing.T) {
	results := tagger.Tag(map[string]string{}, []tagger.Rule{{Tag: "x", Prefixes: []string{"X_"}}})
	if len(results) != 0 {
		t.Errorf("expected no results for empty env")
	}
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
